// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"errors"
	"fmt"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/pangeacyber/pangea-cli/v2/plugins"
	"github.com/spf13/cobra"
)

var PluginSelect = plugins.NewPlugin(selectCmd, []string{"vault", "workspace", "select"})

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select the workspace you want to get secrets from on Pangea Vault",
	Long: `This command selects the workspace you want to link your current directory to a remote directory on Pangea Vault.
You need to do this before you use "pangea vault workspace run" to specify which directory you want to fetch secrets from.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		if workspace == "" {
			return errors.New("empty workspace name not allowed")
		}

		err = selectWorkspace(workspace)
		if err != nil {
			return err
		}

		logger.Printf("Workspace '%s' selected and cached.\n", workspace)
		return nil
	},
}

func init() {
	selectCmd.Flags().StringP("name", "n", "", "folder name on Pangea Vault (Example - /<workspace_name>/dev)")
}

func selectWorkspace(workspace string) error {
	wd := GetWD()

	paths, err := cli.CacheGetPaths()
	if err != nil {
		return err
	}

	if workspace == "" {
		workspace = promptUser("Enter the name of your workspace: ")
	}

	paths[wd] = cli.WorkspaceData{
		Remote: workspace,
	}

	return cli.CacheSetPaths(paths)
}

func promptUser(promptMessage string) string {
	logger.Print(promptMessage)
	var input string
	_, _ = fmt.Scanln(&input)
	return input
}
