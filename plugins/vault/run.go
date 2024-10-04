// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation
package vault

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/spf13/cobra"
)

var PluginRun = plugins.NewPlugin(runCmd, []string{"vault", "workspace", "run"})

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run your application with secrets on Pangea Vault",
	Long: `Run your applications with secrets loaded as environment variables into your application on Pangea Vault.

	For example:
		pangea vault workspace run -- npm run dev
			- will start your node server with secrets loaded in from Pangea`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("no specified command")
		}

		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil {
			log.Fatal(err)
		}

		if workspace == "" {
			workspace = GetWorkspaceFromSettings()
		}

		baseCommand := args[0]
		args = args[1:]
		err = execSubprocess(workspace, baseCommand, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func execSubprocess(workspace string, baseCommand string, args []string) error {
	cmd := exec.Command(baseCommand, args...)
	remoteEnv := GetWorkspaceSecrets(workspace)

	env := make([]string, len(os.Environ())+len(remoteEnv))
	copy(env, os.Environ())
	env = append(env, remoteEnv...)
	cmd.Env = env

	// Set up pipes for stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the subprocess
	err := cmd.Start()
	if err != nil {
		log.Fatal("Error starting subprocess:", err)
		return err
	}

	// Wait for the subprocess to finish
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Error waiting for subprocess:", err)
		return err
	}

	// Subprocess has completed
	return nil
}

func init() {
	runCmd.Flags().StringP("workspace", "w", "", "Overwrite 'workspace' selected with 'pangea vault workspace select'")
}
