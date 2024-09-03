// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v1.2.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Pangea CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
