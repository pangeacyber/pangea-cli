// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	sv "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PluginMigrate = plugins.NewPlugin(migrateCmd, []string{"vault", "workspace", "migrate"})

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate local .env file to Pangea Vault",
	Long: `Migrate your local .env file to Pangea's secure vault.
Simply run "pangea vault workspace migrate -f <path_to_env_file>"

Note: You must select or create a workspace before running migrate.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		envFilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		// TODO: Share viper instance across app
		userConfigViper := viper.New()
		// Initialize Viper
		userConfigViper.SetConfigFile(envFilePath)

		// Enable reading environment variables
		userConfigViper.AutomaticEnv()
		userConfigViper.SetConfigType("dotenv")

		// Read the .env file
		if err := userConfigViper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading .env file: %s", err)
		}

		client, err := CreateVaultService()
		if err != nil {
			return err
		}

		logger.Println("Migrating Secrets 🪄...")

		// Loop through variables from the .env file
		settings := userConfigViper.AllSettings()
		for key, value := range settings {
			v, ok := value.(string)
			if !ok {
				continue
			}
			logger.Println(strings.ToUpper(key))
			_, err := client.SecretStore(
				context.Background(),
				&sv.SecretStoreRequest{
					CommonStoreRequest: sv.CommonStoreRequest{
						Name:   strings.ToUpper(key),
						Folder: workspace,
					},
					Secret: v,
				},
			)
			if err != nil {
				return err
			}
		}

		logger.Printf("Success! All secrets have been migrated to %s in your secure Pangea Vault\n", workspace)
		logger.Println("You can now delete your env file and run `pangea vault workspace run -- $APP_COMMAND`")
		return nil
	},
}

func init() {
	// migrateCmd represents the migrate command
	migrateCmd.Flags().StringP("file", "f", ".env", "env file path (Ex. .env, .env.local)")
	migrateCmd.Flags().StringP("workspace", "w", "", "Overwrite 'workspace' selected with 'pangea vault workspace select'")
}
