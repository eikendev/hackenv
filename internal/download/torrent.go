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

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	progressbar "github.com/schollz/progressbar/v3"
)

const progressInterval = 500 * time.Millisecond

// Torrent downloads the single file described by the metainfo into dir and returns its path.
// The client only leeches: it never uploads, and it stops as soon as the file is complete.
func Torrent(ctx context.Context, metaInfo []byte, dir string, logger *slog.Logger) (string, error) {
	mi, err := metainfo.Load(bytes.NewReader(metaInfo))
	if err != nil {
		return "", fmt.Errorf("cannot parse torrent file: %w", err)
	}

	store := storage.NewFileOpts(storage.NewFileClientOpts{
		ClientBaseDir:   dir,
		PieceCompletion: storage.NewMapPieceCompletion(),
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
		return "", fmt.Errorf("cannot start torrent client: %w", err)
	}
	defer client.Close()

	t, err := client.AddTorrent(mi)
	if err != nil {
		return "", fmt.Errorf("cannot add torrent: %w", err)
	}

	select {
	case <-t.GotInfo():
	case <-ctx.Done():
		return "", ctx.Err()
	}

	if len(t.Files()) != 1 {
		return "", fmt.Errorf("torrent contains %d files, expected exactly one", len(t.Files()))
	}

	t.DownloadAll()

	if err := waitForCompletion(ctx, t); err != nil {
		return "", err
	}

	path := filepath.Join(dir, t.Files()[0].Path())
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("downloaded file missing: %w", err)
	}

	return path, nil
}

func waitForCompletion(ctx context.Context, t *torrent.Torrent) error {
	bar := progressbar.DefaultBytes(t.Length(), "downloading")
	ticker := time.NewTicker(progressInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_ = bar.Set64(t.BytesCompleted())
			if t.BytesMissing() == 0 {
				_ = bar.Finish()
				return nil
			}
		case <-t.Closed():
			return errors.New("torrent was closed before completion")
		}
	}
}
