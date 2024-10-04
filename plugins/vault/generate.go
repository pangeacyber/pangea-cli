package vault

import (
	"os"
	"strconv"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/pangeacyber/pangea-cli-internal/plugins/vault/ed25519"
	"github.com/pangeacyber/pangea-cli-internal/plugins/vault/rsa"
	"github.com/spf13/cobra"
)

var cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generate keys locally",
	Long:  "Generate keys locally",
}

var PluginVaultGenerate = plugins.NewPlugin(cmdGenerate, []string{"vault", "local", "generate"})

func init() {
	cmdGenerate.PersistentFlags().StringP("output", "o", "key.pem", "Output file name to save private key. Public key will be saved on `<filename>.pub`")

	cmdGenerateEd25519 := &cobra.Command{
		Use:   "ed25519",
		Short: "Generate an ED25519 key pair",
		Long:  "Generate an ED25519 key pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			output := cmd.Flag("output").Value.String()
			pub, priv, err := ed25519.GenerateKeyPair()
			if err != nil {
				return err
			}

			b, err := ed25519.EncodePEMPrivateKey(priv)
			if err != nil {
				return err
			}

			err = os.WriteFile(output, b, 0600)
			if err != nil {
				return err
			}

			b, err = ed25519.EncodePEMPublicKey(pub)
			if err != nil {
				return err
			}

			err = os.WriteFile(output+".pub", b, 0600)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmdGenerateRSA := &cobra.Command{
		Use:   "rsa",
		Short: "Generate an RSA key pair",
		Long:  "Generate an RSA key pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			output := cmd.Flag("output").Value.String()

			bz := 4096
			var err error
			fbz := cmd.Flag("bits")
			if fbz != nil {
				bz, err = strconv.Atoi(fbz.Value.String())
				if err != nil {
					return err
				}
			}

			pub, priv, err := rsa.GenerateKeyPair(bz)
			if err != nil {
				return err
			}

			b, err := rsa.EncodePEMPrivateKey(priv)
			if err != nil {
				return err
			}

			err = os.WriteFile(output, b, 0600)
			if err != nil {
				return err
			}

			b, err = rsa.EncodePEMPublicKey(pub)
			if err != nil {
				return err
			}

			err = os.WriteFile(output+".pub", b, 0600)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmdGenerateRSA.Flags().Int("bits", 4096, "Size of the key pair in bits. Possible values: [2048, 3072, 4096].")

	cmdGenerate.AddCommand(
		cmdGenerateEd25519,
		cmdGenerateRSA,
	)
}
