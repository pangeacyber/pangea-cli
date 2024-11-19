package profile

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/spf13/cobra"
)

var profileInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about current selected profile",
	Long:  "Information about current selected profile",
	Run: func(cmd *cobra.Command, args []string) {
		profile, name, err := cli.GetCurrentProfile()
		if err != nil {
			fmt.Printf("Failed to get current profile information. Error[%v]. Select a profile running 'pangea admin profile select' command.\n", err)
			return
		}

		s, err := cli.IndentedString(profile)
		if err != nil {
			fmt.Printf("Failed to print profile info. Error[%v]\n", err)
			return
		}

		fmt.Printf("Current profile: %s\nData:\n", name)
		fmt.Println(s)
	},
}
