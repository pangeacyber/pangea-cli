/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var folderName string
		isPathExists, config, currentDir := utils.CheckPathExists()
		if isPathExists {
			folderName = config.Paths[currentDir].Remote
		} else {
			log.Fatalln("Pangea project is not setup. Run `pange select` to select your project.")
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

		client := utils.CreateVaultAPIClient()

		fmt.Println("Migrating Secrets 🪄...")

		// Loop through variables from the .env file
		settings := userConfigViper.AllSettings()
		for key, value := range settings {
			fmt.Println(strings.ToUpper(key))
			resp, err := client.R().
				SetBody(fmt.Sprintf(`{"name":"%s", "secret":"%s", "folder":"%s"}`, strings.ToUpper(key), value, folderName)).
				Post("https://vault.aws.us.pangea.cloud/v1/secret/store")
			if err != nil {
				log.Fatal(err)
			}

			if resp.IsError() {
				if resp.StatusCode() == 400 {
					err = fmt.Errorf("Error: Secret %s already exists in your project at %s.\nPlease go to your project at https://console.pangea.cloud/service/vault/data to rotate (update) it.", strings.ToUpper(key), folderName)
					fmt.Println(err)
					os.Exit(1)
				} else {
					log.Fatal("Error migrating secrets to your vault. More info:\n", err)
				}
			}
		}

		fmt.Printf("Success! All secrets have been migrated to %s in your secure Pangea Vault", folderName)
		fmt.Println("You can now delete your env file and run `pangea run -c $APP_COMMAND`")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringP("file", "f", ".env", "env file path (Ex. .env, .env.local)")

}
