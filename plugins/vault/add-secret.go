// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	sv "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
	"github.com/spf13/cobra"
)

var PluginAddSecret = plugins.NewPlugin(addSecretCmd, []string{"vault", "workspace", "add-secret"})

var addSecretCmd = &cobra.Command{
	Use:   "add-secret",
	Short: "List all available worspace in Pangea Vault inside <root>",
	Long:  "List all available worspace in Pangea Vault inside <root>",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		value, err := cmd.Flags().GetString("value")
		if err != nil {
			return err
		}

		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil {
			return err
		}

		if workspace == "" {
			workspace = GetWorkspaceFromSettings()
		}

		if workspace == "" {
			return errors.New("workspace not found. Please use 'pangea vault workspace select' to choose the workspace you would like to upload secrets to or set '--workspace' flag")
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()

		vaultClient, err := CreateVaultService()
		if err != nil {
			return err
		}

		_, err = vaultClient.SecretStore(
			ctx,
			&sv.SecretStoreRequest{
				Secret: value,
				CommonStoreRequest: sv.CommonStoreRequest{
					Name:   name,
					Folder: workspace,
				},
			})

		if err != nil {
			return err
		}

		fmt.Printf("Secret added successfully. Name: %s\n", name)
		return nil
	},
}

func init() {
	// addSecretCmd represents the select command
	addSecretCmd.Flags().StringP("value", "v", "", "Secret's value to add to the workspace")
	_ = addSecretCmd.MarkFlagRequired("value")
	addSecretCmd.Flags().StringP("name", "n", "", "Secret's name to add to the workspace")
	_ = addSecretCmd.MarkFlagRequired("name")
	addSecretCmd.Flags().StringP("workspace", "w", "", "Overwrite 'workspace' selected with 'pangea vault workspace select'")
}
