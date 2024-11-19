package sync

import (
	"github.com/pangeacyber/pangea-cli/v2/plugins"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault workspace to external services.",
	Long:  "Sync secrets from Vault workspace to external services.",
}

var PluginSync = plugins.NewPlugin(syncCmd, []string{"sync"})
