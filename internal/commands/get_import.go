package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
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

func (c *GetCommand) verifyImportSource(info *images.DownloadInfo) (os.FileInfo, error) {
	checksum, err := calculateFileChecksum(c.File)
	if err != nil {
		return nil, err
	}

	if checksum != info.Checksum {
		slog.Error("found invalid checksum", "path", c.File, "expected", info.Checksum, "actual", checksum)
		return nil, fmt.Errorf("%s does not match the current %s release (expected %s, actual %s)", c.File, info.Filename, info.Checksum, checksum)
	}

	slog.Info("validated checksum")

	stat, err := os.Stat(c.File)
	if err != nil {
		return nil, fmt.Errorf("cannot stat %s: %w", c.File, err)
	}

	return stat, nil
}

// importImage verifies a locally available image against the signed checksum and copies it into the data directory.
// It reports whether the image was imported.
func (c *GetCommand) importImage(ctx context.Context, info *images.DownloadInfo, localPath string) (bool, error) {
	slog.Info("importing image", "source", c.File, "path", localPath)

	stat, err := c.verifyImportSource(info)
	if err != nil {
		return false, err
	}

	if !c.Yes {
		ok, err := confirm(fmt.Sprintf("Copy %s (%.1f GiB) to %s?", c.File, float64(stat.Size())/(1<<30), localPath))
		if err != nil {
			return false, err
		}
		if !ok {
			slog.Info("skipping image import", "reason", "declined")
			return false, nil
		}
	}

	partial := localPath + partialSuffix

	if err := copyFile(ctx, c.File, partial); err != nil {
		removePartial(partial)
		slog.Error("failed to copy image", "source", c.File, "path", partial, "err", err)
		return false, err
	}

	if err := os.Rename(partial, localPath); err != nil {
		removePartial(partial)
		slog.Error("failed to move image into place", "from", partial, "to", localPath, "err", err)
		return false, fmt.Errorf("cannot move %s to %s: %w", partial, localPath, err)
	}

	slog.Info("imported image", "path", localPath)
	return true, nil
}
