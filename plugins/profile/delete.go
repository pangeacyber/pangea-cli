package profile

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli-internal/cli"
	"github.com/spf13/cobra"
)

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete profile",
	Long:  "Delete profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName, _ := cmd.Flags().GetString("name")

		ok := askConfirmation(profileName)
		if !ok {
			fmt.Printf("Delete of profile '%s' aborted\n", profileName)
			return
		}

		err := cli.DeleteProfile(profileName)
		if err != nil {
			fmt.Printf("Failed to delete profile. Error[%v]\n", err)
			return
		}

		fmt.Printf("Profile '%s' delete successfully.\n", profileName)
	},
}

func init() {
	profileDeleteCmd.Flags().StringP("name", "n", "", "Profile name")
}

func askConfirmation(name string) bool {
	fmt.Printf("To confirm deletion of profile '%s' write its name an press enter:\n", name)
	in := cli.ReadStdin()
	return in == name
}
