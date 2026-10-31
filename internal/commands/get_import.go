package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	progressbar "github.com/schollz/progressbar/v3"

	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
)

func confirm(question string) (bool, error) {
	fmt.Fprintf(os.Stderr, "%s [y/N] ", question)

	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("cannot read answer: %w", err)
	}

	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func copyFile(ctx context.Context, src, dst string) error {
	in, err := os.Open(src) //#nosec G304
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", src, err)
	}
	defer handling.Close(in)

	stat, err := in.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat %s: %w", src, err)
	}

	out, err := os.Create(dst) //#nosec G304
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", dst, err)
	}
	defer handling.Close(out)

	bar := progressbar.DefaultBytes(stat.Size(), "copying")

	if _, err := io.Copy(io.MultiWriter(out, bar), &contextReader{ctx: ctx, r: in}); err != nil {
		return fmt.Errorf("cannot copy %s to %s: %w", src, dst, err)
	}

	return nil
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (c *contextReader) Read(p []byte) (int, error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	return c.r.Read(p)
}

func (c *GetCommand) importTarget(imageType string) (string, error) {
	image, err := images.GetImageDetails(imageType)
	if err != nil {
		slog.Error("failed to resolve image", "type", imageType, "err", err)
		return "", fmt.Errorf("cannot resolve image details for %q: %w", imageType, err)
	}

	version := image.FileVersion(filepath.Base(c.File))
	if version == "" {
		slog.Error("missing version in filename", "path", c.File)
		return "", fmt.Errorf("cannot determine release of %s from its name; keep the vendor's filename", c.File)
	}

	localPath, err := image.GetLocalPath(version)
	if err != nil {
		slog.Error("failed to resolve image path", "image", image.DisplayName, "version", version, "err", err)
		return "", fmt.Errorf("cannot resolve local path for %s %s: %w", image.DisplayName, version, err)
	}

	return localPath, nil
}

func (c *GetCommand) shouldSkipImport(localPath string) (bool, error) {
	_, err := os.Stat(localPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		slog.Error("failed to stat image", "path", localPath, "err", err)
		return false, fmt.Errorf("cannot stat %s: %w", localPath, err)
	}

	if !c.Force {
		slog.Info("image already exists", "path", localPath, "hint", "use --force to overwrite")
		return true, nil
	}

	return false, nil
}

// importFile copies a locally available image into the data directory without verifying it.
func (c *GetCommand) importFile(ctx context.Context, imageType string) error {
	localPath, err := c.importTarget(imageType)
	if err != nil {
		return err
	}

	skip, err := c.shouldSkipImport(localPath)
	if err != nil || skip {
		return err
	}

	slog.Warn("skipping checksum verification", "reason", "imported image", "path", c.File)
	slog.Info("importing image", "source", c.File, "path", localPath)

	stat, err := os.Stat(c.File)
	if err != nil {
		return fmt.Errorf("cannot stat %s: %w", c.File, err)
	}

	if !c.Yes {
		ok, err := confirm(fmt.Sprintf("Copy %s (%.1f GiB) to %s?", c.File, float64(stat.Size())/(1<<30), localPath))
		if err != nil {
			return err
		}
		if !ok {
			slog.Info("skipping image import", "reason", "declined")
			return nil
		}
	}

	partial := localPath + partialSuffix

	if err := copyFile(ctx, c.File, partial); err != nil {
		removePartial(partial)
		slog.Error("failed to copy image", "source", c.File, "path", partial, "err", err)
		return err
	}

	if err := os.Rename(partial, localPath); err != nil {
		removePartial(partial)
		slog.Error("failed to move image into place", "from", partial, "to", localPath, "err", err)
		return fmt.Errorf("cannot move %s to %s: %w", partial, localPath, err)
	}

	slog.Info("imported image", "path", localPath)
	hintSELinuxLabels()
	return nil
}
