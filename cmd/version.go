/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Pangea CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pangea v1.0 Beta")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
