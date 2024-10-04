package plugins

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type plugin struct {
	cmd               *cobra.Command
	cmdPath           []string
	serviceAssociated string
}

type CommandProvider interface {
	GetCommand([]string) *cobra.Command
}

func NewPlugin(c *cobra.Command, p []string, options ...PluginOption) *plugin {
	return &plugin{
		cmd:     c,
		cmdPath: p,
	}
}

type PluginOption func(*plugin) error

// Use this to associate a plugin with a service. So if failed to load service, this plugin won't load.
// A plugin is considerated Service associated if get a service command and modify its behaviour
func AssociateService(service string) PluginOption {
	return func(p *plugin) error {
		p.serviceAssociated = service
		return nil
	}
}

func (p *plugin) GetCommand() *cobra.Command {
	return p.cmd
}

func (p *plugin) GetCommandPath() []string {
	return p.cmdPath
}

func (p *plugin) Init(b CommandProvider) error {
	return nil
}

func (p *plugin) GetServiceAssociated() string {
	return p.serviceAssociated
}

// Plugin interface represents the functionality that plugins can provide.
type Plugin interface {
	GetCommand() *cobra.Command
	GetCommandPath() []string
	Init(CommandProvider) error
	GetServiceAssociated() string
}

func CommandPathToString(path []string) string {
	cmd := ""
	for _, p := range path {
		cmd = fmt.Sprintf("%s%s ", cmd, p)
	}
	return strings.TrimSpace(cmd)
}
