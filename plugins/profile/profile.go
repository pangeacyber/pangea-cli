package profile

import (
	"github.com/pangeacyber/pangea-cli/v2/plugins"
	"github.com/spf13/cobra"
)

var PluginProfile = plugins.NewPlugin(profileCmd, []string{"admin", "profile"})

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "List of 'profile' commands",
	Long:  "List of 'profile' commands",
}

func init() {
	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileSelectCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileInfoCmd)
	profileCmd.AddCommand(profileUpdateCmd)
	profileCmd.AddCommand(profileDeleteCmd)
}
