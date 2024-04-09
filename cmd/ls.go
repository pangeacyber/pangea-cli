// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "View all secrets for the selected workspace",
	Long: `pangea ls
	
	shows you all the secrets in your current selected workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		remoteEnv := Get_env()

		for _, envVar := range remoteEnv {
			fmt.Println(envVar)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
