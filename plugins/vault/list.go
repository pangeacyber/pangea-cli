// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	sv "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
	"github.com/spf13/cobra"
)

var PluginList = plugins.NewPlugin(listCmd, []string{"vault", "workspace", "list"})

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available secrets workspaces",
	Long:  "List worspaces availables in Pangea Vault inside the folder set in 'root' flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			return err
		}

		details, err := cmd.Flags().GetBool("details")
		if err != nil {
			return err
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()

		filter := map[string]string{
			"type":   "folder",
			"folder": root,
		}

		vaultClient, err := CreateVaultService()
		if err != nil {
			return err
		}

		resp, err := vaultClient.List(
			ctx,
			&sv.ListRequest{
				Filter: filter,
			})

		if err != nil {
			return err
		}

		for _, item := range resp.Result.Items {
			fmt.Println(item.Name)
			if details {
				fmt.Printf("\tid: %s  created_at: %s\n", item.ID, item.CreatedAt)
			}
		}

		return nil
	},
}

func init() {
	// selectCmd represents the select command
	listCmd.Flags().StringP("root", "r", "/", "Root folder name on Pangea Vault to list available workspaces inside it (Example - <root>/<workspace_name>)")
	listCmd.Flags().BoolP("details", "d", false, "Print extra information about each workspace.")
}
