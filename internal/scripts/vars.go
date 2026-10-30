// Package scripts contains scripts that can be run from within the tool.
package scripts

import (
	_ "embed"
)

// ApplyLabelsScript is a script that applies necessary SELinux labels to images and shared folder.
//
//go:embed bin/hackenv_applylabels
var ApplyLabelsScript string
