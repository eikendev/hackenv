package paths

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/eikendev/hackenv/internal/constants"
)

func GetDataFilePath(file string) string {
	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, file))
	if err != nil {
		log.Fatalf("Cannot access data directory: %s\n", err)
	}

	return path
}

func EnsureDirExists(path string) {
	err := os.MkdirAll(path, 0660)
	if err != nil {
		log.Fatalf("Cannot create directory: %s\n", err)
	}
}
