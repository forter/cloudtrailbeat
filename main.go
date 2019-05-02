package main

import (
	"os"

	"github.com/forter/cloudtrailbeat/cmd"

	_ "github.com/forter/cloudtrailbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
