// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"errors"
	"fmt"
	"os"

	"github.com/pangeacyber/pangea-cli/v2/plugins"
	"github.com/spf13/cobra"
)

var PluginStoreFromFile = &pluginStoreFromFile{
	serviceAssociated: "vault",
}

type pluginStoreFromFile struct {
	pluginCmd         *cobra.Command
	serviceAssociated string
}

func (p *pluginStoreFromFile) GetCommandPath() []string {
	return []string{"vault", "v1", "/key/store", "file"}
}

func (p *pluginStoreFromFile) GetCommand() *cobra.Command {
	return p.pluginCmd
}

func (p *pluginStoreFromFile) Init(b plugins.CommandProvider) error {
	cmd := b.GetCommand([]string{"vault", "v1", "/key/store"})
	if cmd == nil {
		return errors.New("not available command")
	}

	cmd.Use = "file"
	cmd.Short = "Store key in Vault from private and public key PEM file"
	cmd.Long = "Store key(s) in Vault from PEM file"
	cmd.Args = cobra.ExactArgs(1)
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		path := ""
		for _, p := range p.GetCommandPath() {
			path = fmt.Sprintf("%s%s ", path, p)
		}
		logger.Println("Usage:")
		logger.Printf("    pangea %s<pem-file-name> [flags]\n", path)
		logger.Println("\nFlags:")
		logger.Println(cmd.Flags().FlagUsages())
		return nil
	})
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		tflag := cmd.Flag("type")
		if tflag == nil {
			return errors.New("`type` need to be set")
		}

		switch tflag.Value.String() {
		case "fpe":
		case "symmetric_key":
			// If us is symmetric or fpe
			fkey := cmd.Flag("key")
			if fkey == nil {
				return errors.New("no `key` flag available")
			}
			err = fkey.Value.Set(string(content))
			if err != nil {
				return err
			}
		case "asymmetric_key":
			// Should set public key too?
			fkey := cmd.Flag("private_key")
			if fkey == nil {
				return errors.New("no `private_key` flag available")
			}
			err = fkey.Value.Set(string(content))
			if err != nil {
				return err
			}
			fkey.Changed = true

			fkey = cmd.Flag("public_key")
			if fkey == nil {
				return errors.New("no `public_key` flag available")
			}

			content, err := os.ReadFile(args[0] + ".pub")
			if err != nil {
				return err
			}
			err = fkey.Value.Set(string(content))
			if err != nil {
				return err
			}
			fkey.Changed = true
		default:
			return fmt.Errorf("not supported `type`: %s", tflag.Value)
		}

		return nil
	}

	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	p.pluginCmd = cmd
	return nil
}

func (p *pluginStoreFromFile) GetServiceAssociated() string {
	return p.serviceAssociated
}
