package main

import (
	"os"

	"github.com/TranscriptionFactory/shift-log/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
