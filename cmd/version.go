// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Pangea CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.1.1 Beta")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
