package loader

import (
	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/pangeacyber/pangea-cli-internal/plugins/intel"
	"github.com/pangeacyber/pangea-cli-internal/plugins/profile"
	"github.com/pangeacyber/pangea-cli-internal/plugins/sync"
	"github.com/pangeacyber/pangea-cli-internal/plugins/updates"
	"github.com/pangeacyber/pangea-cli-internal/plugins/utils"
	"github.com/pangeacyber/pangea-cli-internal/plugins/vault"
)

func LoadPlugins() []plugins.Plugin {
	return []plugins.Plugin{
		// Add yours plugins to this list
		utils.PluginBase64,
		vault.PluginMigrate,
		vault.PluginRun,
		vault.PluginSelect,
		vault.PluginCreate,
		vault.PluginListSecrets,
		vault.PluginStoreFromFile,
		vault.PluginVaultGenerate,
		intel.PluginIntelFilePatternReputation,
		sync.PluginSync,
		sync.PluginSyncVercel,
		vault.PluginList,
		vault.PluginAddSecret,
		vault.PluginWorkspace,
		updates.PluginCheckUpdate,
		updates.PluginUpdate,
		profile.PluginProfile,
	}
}
