package profile

import (
	"fmt"

	"github.com/pangeacyber/pangea-cli-internal/cli"
	"github.com/spf13/cobra"
)

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	Long:  "List available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := cli.ListProfiles()
		if err != nil {
			fmt.Printf("Unable to list profiles. Error[%v]\n", err)
			return
		}

		if len(profiles) == 0 {
			fmt.Println("No available profiles. Create one running 'pangea admin profile create' command.")
			return
		}

		fmt.Println("Available profiles:")
		for _, p := range profiles {
			fmt.Printf("\t%s\n", p)
		}
	},
}
