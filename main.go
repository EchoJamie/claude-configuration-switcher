package main

import (
	"os"

	"github.com/EchoJamie/ccs/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
