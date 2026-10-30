// Package download fetches large files through BitTorrent.
package download

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	g "github.com/anacrolix/generics"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	progressbar "github.com/schollz/progressbar/v3"
)

const progressInterval = 500 * time.Millisecond

// Torrent downloads the single file described by the metainfo into dir and returns its path.
// The client only leeches: it never uploads, and it stops as soon as the file is complete.
// Progress is persisted in dir, so a repeated call with the same dir resumes the download.
func Torrent(ctx context.Context, metaInfo []byte, dir string, logger *slog.Logger) (string, error) {
	mi, err := metainfo.Load(bytes.NewReader(metaInfo))
	if err != nil {
		return "", fmt.Errorf("cannot parse torrent file: %w", err)
	}

	length, err := fetch(ctx, mi, dir, logger)
	if err != nil {
		return "", err
	}

	return findDownloadedFile(dir, length)
}

func fetch(ctx context.Context, mi *metainfo.MetaInfo, dir string, logger *slog.Logger) (int64, error) {
	store := storage.NewFileOpts(storage.NewFileClientOpts{
		ClientBaseDir: dir,
		UsePartFiles:  g.Some(false),
	})
	defer func() { _ = store.Close() }()

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = dir
	cfg.DefaultStorage = store
	cfg.ListenPort = 0
	cfg.NoUpload = true
	cfg.Seed = false
	cfg.NoDefaultPortForwarding = true
	cfg.Slogger = logger

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return 0, fmt.Errorf("cannot start torrent client: %w", err)
	}
	defer client.Close()

	t, err := client.AddTorrent(mi)
	if err != nil {
		return 0, fmt.Errorf("cannot add torrent: %w", err)
	}

	select {
	case <-t.GotInfo():
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	if len(t.Files()) != 1 {
		return 0, fmt.Errorf("torrent contains %d files, expected exactly one", len(t.Files()))
	}

	t.DownloadAll()

	if err := waitForCompletion(ctx, t); err != nil {
		return 0, err
	}

	return t.Length(), nil
}

func waitForCompletion(ctx context.Context, t *torrent.Torrent) error {
	bar := progressbar.DefaultBytes(t.Length(), "downloading")
	ticker := time.NewTicker(progressInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.Complete().On():
			_ = bar.Finish()
			return nil
		case <-ticker.C:
			_ = bar.Set64(t.BytesCompleted())
		case <-t.Closed():
			return errors.New("torrent was closed before completion")
		}
	}
}

// findDownloadedFile returns the regular file in dir whose size matches the torrent length.
func findDownloadedFile(dir string, length int64) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("cannot list %s: %w", dir, err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() || info.Size() != length {
			continue
		}
		return filepath.Join(dir, entry.Name()), nil
	}

	return "", fmt.Errorf("no file of %d bytes found in %s", length, dir)
}
