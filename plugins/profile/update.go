package profile

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli-internal/cli"
	"github.com/spf13/cobra"
)

var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a profile token or domain",
	Long:  "Update a profile token or domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName, _ := cmd.Flags().GetString("profile")
		serviceName, _ := cmd.Flags().GetString("service")
		update := false

		token, _ := cmd.Flags().GetString("token")
		if token != "" {
			err := cli.SaveToken(profileName, serviceName, token)
			if err != nil {
				return err
			}
			update = true
		}

		domain, _ := cmd.Flags().GetString("domain")
		if domain != "" {
			err := cli.SaveDomain(profileName, serviceName, domain)
			if err != nil {
				return err
			}
			update = true
		}

		if update && profileName == "" {
			p, err := cli.GetCurrentProfileName()
			if err == nil {
				profileName = p
			}
		}

		if update {
			fmt.Printf("Profile '%s' updated successfully.\n", profileName)
		}

		if !update {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	profileUpdateCmd.Flags().StringP("profile", "a", "", "Profile name to be updated. If omitted current selected profile will be updated.")
	profileUpdateCmd.Flags().StringP("service", "s", "default", "Service name to be updated. If omitted 'default' service will be updated.")
	profileUpdateCmd.Flags().StringP("token", "t", "", "Token to be saved")
	profileUpdateCmd.Flags().StringP("domain", "d", "", "Domain to be saved")
}
