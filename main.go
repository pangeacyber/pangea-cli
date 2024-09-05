// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package main

import (
	"fmt"
	"os"

	"github.com/pangeacyber/pangea-cli/cmd"
	"github.com/pangeacyber/pangea-cli/updates"
)

func main() {
	_, _, err := updates.CheckAvailableVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error checking new version available.")
	}
	cmd.Execute()
}
