package cmd

import (
	"os"

	"github.com/forter/cloudtrailbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/libbeat/cmd/instance"
)

// Name of this beat
var Name = "cloudtrailbeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name})

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
