package main

import (
	"github.com/denizgursoy/gotouch/internal/commands"
)

var (
	Version     = "v0.0.0"
	BuildCommit = "0000000"
	BuildDate   = "0000-00-00T00:00:00Z"
)

func main() {
	commands.Execute(commands.BuildInfo{Version: Version, BuildCommit: BuildCommit, BuildDate: BuildDate})
}
