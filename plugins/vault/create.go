// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package vault

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-cli/v2/plugins"
	sv "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
	"github.com/spf13/cobra"
)

var PluginCreate = plugins.NewPlugin(createCmd, []string{"vault", "workspace", "create"})

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Pangea secrets workspace",
	Long:  `Creates a workspace in your pangea vault which let's you store your secrets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		if workspaceName == "" {
			logger.Print("Enter the name of your workspace: ")
			_, err := fmt.Scanln(&workspaceName)
			if err != nil {
				return err
			}
		}

		client, err := CreateVaultService()
		if err != nil {
			return err
		}

		_, err = client.FolderCreate(
			context.Background(),
			&sv.FolderCreateRequest{
				Name:   workspaceName,
				Folder: "",
			})
		if err != nil {
			return err
		}

		err = selectWorkspace(fmt.Sprintf("/%s/", workspaceName))
		if err != nil {
			logger.Fatalf("Failed to select workspace '%s'. Error[%v]\n", workspaceName, err)
		}
		logger.Printf("workspace created at %s in Pangea vault\n\nRun `pangea vault workspace migrate -f .env` to migrate your .env file to your workspace", fmt.Sprintf("/%s/", workspaceName))
		return nil
	},
}

func init() {
	// createCmd represents the create command
	createCmd.Flags().StringP("name", "n", "", "workspace Name (used as folder name)")
}
