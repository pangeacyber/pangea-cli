/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

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

		envFilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatalln(err)
		}

		// Initialize Viper
		viper.SetConfigFile(envFilePath)

		// Enable reading environment variables
		viper.AutomaticEnv()

		// Read the .env file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading .env file: %s\n", err)
			os.Exit(1)
		}

		client := utils.CreateVaultAPIClient()

		fmt.Println("Migrating Secrets 🪄...")

		// TODO: Make the folder dynamically read from cache_path.
		folderName := "/secrets/bryan/dev/"

		// Loop through variables from the .env file
		settings := viper.AllSettings()
		for key, value := range settings {
			_, err := client.R().
				SetBody(fmt.Sprintf(`{"name":"%s", "secret":"%s", "folder":"%s"}`, key, value, folderName)).
				Post("https://vault.aws.us.pangea.cloud/v1/secret/store")
			if err != nil {
				log.Fatal(err)
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
