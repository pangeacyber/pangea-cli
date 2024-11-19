package utils

import (
	"encoding/base64"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/pangeacyber/pangea-cli/v2/plugins"
	"github.com/spf13/cobra"
)

var PluginBase64 = plugins.NewPlugin(getBase64cmd(), []string{"utils", "base64"})
var logger = cli.GetLogger()

func getBase64cmd() *cobra.Command {
	base64cmd := &cobra.Command{
		Use:   "base64",
		Short: "Base64 functions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	encode := &cobra.Command{
		Use:   "encode",
		Short: "encode string in base64",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Println(base64.StdEncoding.EncodeToString([]byte(args[0])))
			return nil
		},
	}

	decode := &cobra.Command{
		Use:   "decode",
		Short: "decode a base64 string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			logger.Println(string(d))
			return nil
		},
	}

	// Add commands
	base64cmd.AddCommand(encode, decode)
	return base64cmd
}
