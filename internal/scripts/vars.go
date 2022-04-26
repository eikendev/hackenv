package scripts

import (
	_ "embed"
)

//go:embed bin/hackenv_createbridge
var CreateBridgeScript string

//go:embed bin/hackenv_removebridge
var RemoveBridgeScript string

//go:embed bin/hackenv_applylabels
var ApplyLabelsScript string
