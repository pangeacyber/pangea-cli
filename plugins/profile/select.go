package profile

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/spf13/cobra"
)

var profileSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select an existing profile",
	Long:  "Select an existing profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName, _ := cmd.Flags().GetString("profile")
		err := cli.SelectProfile(profileName)
		if err != nil {
			fmt.Printf("Failed to select profile. Error[%v]. Create a profile running 'pangea admin profile create' command\n", err)
			return
		}

		fmt.Printf("Profile '%s' selected successfully.", profileName)
	},
}

func init() {
	profileSelectCmd.Flags().StringP("profile", "p", "", "Associate token to 'profile'")
}
