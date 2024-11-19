package profile

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/spf13/cobra"
)

var profileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new profile",
	Long:  "Create a new profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName, _ := cmd.Flags().GetString("name")
		err := cli.CreateProfile(profileName)
		if err != nil {
			fmt.Printf("Failed to create profile. Error[%v]\n", err)
			return
		}

		fmt.Printf("Profile '%s' created successfully.\n", profileName)
	},
}

func init() {
	profileCreateCmd.Flags().StringP("name", "n", "", "Profile name")
}
