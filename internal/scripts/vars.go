// Package scripts contains scripts that can be run from within the tool.
package scripts

import (
	_ "embed"
)

// CreateBridgeScript is a script that creates the bridge for network communication.
//
//go:embed bin/hackenv_createbridge
var CreateBridgeScript string

// RemoveBridgeScript is a script that removes the bridge for network communication.
//
//go:embed bin/hackenv_removebridge
var RemoveBridgeScript string

// ApplyLabelsScript is a script that applies necessary SELinux labels to images and shared folder.
//
//go:embed bin/hackenv_applylabels
var ApplyLabelsScript string
