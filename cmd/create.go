/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pangea secrets workspace",
	Long:  `Creates a workspace in your pangea vault which let's you store your secrets.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("workspace-name", "n", "", "workspace Name (used as folder name)")

	createCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return validateInput(cmd.Flags())
	}

}

func validateInput(flags *pflag.FlagSet) error {
	workspaceName, err := flags.GetString("workspace-name")
	if err != nil {
		return err
	}

	if workspaceName == "" {
		fmt.Print("Enter the name of your workspace: ")
		_, err := fmt.Scanln(&workspaceName)
		if err != nil {
			return err
		}
	}

	client, pangeaDomain := utils.CreateVaultAPIClient()
	_, err = client.R().
		SetBody(fmt.Sprintf(`{"name":"%s", "folder":"%s"}`, workspaceName, "/secrets/")).
		Post(fmt.Sprintf("https://vault.%s/v1/folder/create", pangeaDomain))

	if err != nil {
		log.Fatalln(err)
	}

	SelectWorkspace(fmt.Sprintf("/secrets/%s/", workspaceName))

	fmt.Printf("workspace created at %s in Pangea vault\n\nRun `pangea migrate -f .env` to migrate your .env file to your workspace", fmt.Sprintf("/secrets/%s/", workspaceName))

	return nil
}
