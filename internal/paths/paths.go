package paths

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"

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
	err := os.MkdirAll(path, 0o660)
	if err != nil {
		log.Fatalf("Cannot create directory: %s\n", err)
	}
}

func EnsurePostbootExists(path string) bool {
	path = fmt.Sprintf("%s/postboot.sh", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Infof("%s doesn't exists", path)
		return false
	}
	return true
}

func GetCmdPathOrExit(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatalf("Command not found: %s\n", cmd)
	}

	return path
}

func GetCmdPath(cmd string) (string, error) {
	return exec.LookPath(cmd)
}
