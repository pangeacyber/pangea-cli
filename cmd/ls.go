/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "View all secrets for the selected projects",
	Long: `pangea ls
	
	shows you all the secrets in your current selected project.`,
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
