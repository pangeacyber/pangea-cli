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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().Bool("create-envs", false, "Create dev, staging, and prod environments")
	createCmd.Flags().StringP("project-name", "n", "", "Project Name (used as folder name)")

	createCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return validateInput(cmd.Flags())
	}

}

func validateInput(flags *pflag.FlagSet) error {
	shouldCreateEnvironments, err := flags.GetBool("create-envs")
	if err != nil {
		return err
	}

	fmt.Println(shouldCreateEnvironments)

	projectName, err := flags.GetString("project-name")

	if projectName == "" {
		fmt.Print("Enter the name of your project: ")
		_, err := fmt.Scanln(&projectName)
		if err != nil {
			return err
		}
	}

	if !shouldCreateEnvironments {
		fmt.Print("[Recommended] Do you want to create dev, staging, and prod environments (y/n): ")
		var confirm string
		_, err := fmt.Scanln(&confirm)
		if err != nil {
			return err
		}

		if confirm != "y" && confirm != "n" {
			return fmt.Errorf("invalid input. Only 'y' or 'n' are allowed")
		}

		client := utils.CreateVaultAPIClient()

		if confirm == "y" {

			envs := []string{"dev", "stg", "prd"}

			for _, env := range envs {

				_, err := client.R().
					SetBody(fmt.Sprintf(`{"name":"%s", "folder":"%s"}`, env, fmt.Sprintf("secrets/%s/", projectName))).
					Post("https://vault.aws.us.pangea.cloud/v1/folder/create")
				if err != nil {
					log.Fatalln(err)
				}
			}

			fmt.Printf("Created dev, stg, and prd environments at %s in your Pangea Vault", fmt.Sprintf("/secrets/%s/", projectName))
		} else {
			_, err := client.R().
				SetBody(fmt.Sprintf(`{"name":"%s", "folder":"%s"}`, "default", fmt.Sprintf("secrets/%s/", projectName))).
				Post("https://vault.aws.us.pangea.cloud/v1/folder/create")

			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Created a default environment at %s in your Pangea Vault", fmt.Sprintf("/secrets/%s/", projectName))
		}
	}

	return nil
}
