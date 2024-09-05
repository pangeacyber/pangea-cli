// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "View all secrets for the selected workspace",
	Long: `pangea ls
	
	shows you all the secrets in your current selected workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		show, err := cmd.Flags().GetBool("show-secrets")
		if err != nil {
			logger.Fatal(err)
		}

		remoteEnv := GetEnv()

		for _, envVar := range remoteEnv {
			if !show {
				envVar = strings.Split(envVar, "=")[0] + "=********"
			}
			logger.Println(envVar)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("show-secrets", "s", false, "Show the secret values")

}
