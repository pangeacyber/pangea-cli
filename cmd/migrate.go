// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate local .env file to Pangea Vault",
	Long: `Migrate your local .env file to Pangea's secure vault.
	Simply run "pangea run -f <path_to_env_file>"
	
	Note: You must select or create a workspace before running migrate.`,
	Run: func(cmd *cobra.Command, args []string) {

		var folderName string
		isPathExists, config, currentDir := utils.CheckPathExists()
		if isPathExists {
			folderName = config.Paths[currentDir].Remote
		} else {
			fmt.Println("Pangea workspace is not setup. Run `pange select` to select your workspace.")
			// Error exit
			os.Exit(0)
		}

		envFilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatalln(err)
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
			fmt.Printf("Error reading .env file: %s\n", err)
			os.Exit(1)
		}

		client, pangeaDomain := utils.CreateVaultAPIClient()

		fmt.Println("Migrating Secrets 🪄...")

		// Loop through variables from the .env file
		settings := userConfigViper.AllSettings()
		for key, value := range settings {
			fmt.Println(strings.ToUpper(key))
			resp, err := client.R().
				SetBody(fmt.Sprintf(`{"name":"%s", "secret":"%s", "folder":"%s"}`, strings.ToUpper(key), value, folderName)).
				Post(fmt.Sprintf("https://vault.%s/v1/secret/store", pangeaDomain))
			if err != nil {
				log.Fatal(err)
			}

			if resp.IsError() {
				if resp.StatusCode() == 400 {
					err = fmt.Errorf("Error: Secret %s already exists in your workspace at %s.\nPlease go to your workspace at https://console.pangea.cloud/service/vault/data to rotate (update) it.", strings.ToUpper(key), folderName)
					fmt.Println(err)
					// Error exit
					os.Exit(1)
				} else {
					log.Fatal("Error migrating secrets to your vault. More info:\n", resp.Status())
				}
			}
		}

		fmt.Printf("Success! All secrets have been migrated to %s in your secure Pangea Vault\n", folderName)
		fmt.Println("You can now delete your env file and run `pangea run -c $APP_COMMAND`")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringP("file", "f", ".env", "env file path (Ex. .env, .env.local)")

}
