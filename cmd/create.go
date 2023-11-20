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
	Short: "Create a pangea secrets project",
	Long:  `Creates a project in your pangea vault which let's you store your secrets.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("project-name", "n", "", "Project Name (used as folder name)")

	createCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return validateInput(cmd.Flags())
	}

}

func validateInput(flags *pflag.FlagSet) error {
	projectName, err := flags.GetString("project-name")
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Print("Enter the name of your project: ")
		_, err := fmt.Scanln(&projectName)
		if err != nil {
			return err
		}
	}

	client, pangeaDomain := utils.CreateVaultAPIClient()
	_, err = client.R().
		SetBody(fmt.Sprintf(`{"name":"%s", "folder":"%s"}`, projectName, "/secrets/")).
		Post(fmt.Sprintf("https://vault.%s/v1/folder/create", pangeaDomain))

	if err != nil {
		log.Fatalln(err)
	}

	SelectProject(fmt.Sprintf("/secrets/%s/", projectName))

	fmt.Printf("Project created at %s in Pangea vault\n\nRun `pangea migrate -f .env` to migrate your .env file to your project", fmt.Sprintf("/secrets/%s/", projectName))

	return nil
}
