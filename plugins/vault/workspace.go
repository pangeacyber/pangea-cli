// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/spf13/cobra"
)

var PluginWorkspace = plugins.NewPlugin(workspaceCmd, []string{"vault", "workspace"})

// runCmd represents the run command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Print current workspace name",
	Long:  "Print current workspace name",

	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := GetWorkspaceFromSettings()
		if workspace == "" {
			fmt.Printf("Workspace not found. Please use 'pangea vault workspace select' to choose the workspace you would like to work with or set 'PANGEA_DEFAULT_FOLDER' environment variable\n\n")
		} else {
			fmt.Printf("Current workspace: %s\n\n", workspace)
		}
		return cmd.Help()
	},
}
