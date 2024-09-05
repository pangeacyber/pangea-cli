package cmd

import "github.com/spf13/cobra"

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync environment variables",
	Long:  `Sync environment variables to external services.`,
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
