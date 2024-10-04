// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"log"
	"strings"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/spf13/cobra"
)

var PluginListSecrets = plugins.NewPlugin(lsCmd, []string{"vault", "workspace", "list-secrets"})

var lsCmd = &cobra.Command{
	Use:   "list-secrets",
	Short: "List secrets in selected workspace",
	Long:  "List secrets in selected workspace",
	Run: func(cmd *cobra.Command, args []string) {
		show, err := cmd.Flags().GetBool("show-secrets")
		if err != nil {
			log.Fatal(err)
		}

		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil {
			log.Fatal(err)
		}

		if workspace == "" {
			workspace = GetWorkspaceFromSettings()
		}

		remoteEnv := GetWorkspaceSecrets(workspace)

		for _, envVar := range remoteEnv {
			if !show {
				envVar = strings.Split(envVar, "=")[0] + "=********"
			}
			logger.Println(envVar)
		}
	},
}

func init() {
	// lsCmd represents the ls command
	lsCmd.Flags().BoolP("show-secrets", "s", false, "Show the secret values")
	lsCmd.Flags().StringP("workspace", "w", "", "Overwrite 'workspace' selected with 'pangea vault workspace select'")
}
