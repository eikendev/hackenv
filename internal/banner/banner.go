package banner

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

const banner = `
    __               __
   / /_  ____ ______/ /_____  ____ _   __
  / __ \/ __ ` + "`" + `/ ___/ //_/ _ \/ __ \ | / /
 / / / / /_/ / /__/ ,< /  __/ / / / |/ /
/_/ /_/\__,_/\___/_/|_|\___/_/ /_/|___/
`

func PrintBanner() {
	fmt.Fprint(os.Stderr, banner)
	fmt.Fprintln(os.Stderr, "")

	_, err := color.New(color.FgBlue).Fprintln(os.Stderr, "                @eikendev")
	if err != nil {
		log.Warnf("Cannot print banner: %s", err)
	}

	fmt.Fprintln(os.Stderr, "")
}
