/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
)

type VaultListResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	Result    struct {
		Count int `json:"count"`
		Items []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"result"`
}

type VaultSecretResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	Result    struct {
		CurrentVersion struct {
			Secret  string `json:"secret"`
			State   string `json:"state"`
			Version int    `json:"version"`
		} `json:"current_version"`
		ID        string     `json:"id"`
		ItemState string     `json:"item_state"`
		Type      string     `json:"type"`
		Versions  []struct{} `json:"versions"`
	} `json:"result"`
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run your application with secrets on Pangea Vault",
	Long: `Run your applications with secrets loaded as environment variables into your application on Pangea Vault.
	
	For example:
		pangea run -- npm run dev
			- will start your node server with secrets loaded in from Pangea`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("No specified command")
		}

		baseCommand := args[0]
		args = args[1:]
		err := exec_subprocess(baseCommand, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Get_env() []string {

	var remoteEnv []string

	defaultWorkspacePathExists := os.Getenv("PANGEA_DEFAULT_FOLDER")

	isPathExists, config, currentDir := utils.CheckPathExists()
	if isPathExists || defaultWorkspacePathExists != "" {
		var folderName string
		if defaultWorkspacePathExists != "" {
			folderName = defaultWorkspacePathExists
		} else {
			folderName = config.Paths[currentDir].Remote
		}

		fmt.Printf("Fetching secrets from: %s\n", folderName)

		client, pangeaDomain := utils.CreateVaultAPIClient()

		resp, err := client.R().
			SetBody(fmt.Sprintf(`{"filter": {"folder":"%s"}}`, folderName)).
			Post(fmt.Sprintf("https://vault.%s/v1/list", pangeaDomain))
		if err != nil {
			log.Fatalln(err)
		}

		var response VaultListResponse
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			log.Fatal("Error fetching secrets from Pangea")
		}
		if response.Status == "Unauthorized" {
			log.Fatal("Unauthorized! Please run `pangea login` to get a new token.")
		}

		// Create a list of secret type IDs
		secretIDs := make(map[string]interface{})
		for _, item := range response.Result.Items {
			if item.Type == "secret" {
				secretIDs[item.ID] = item.Name
			}
		}

		// Print the list of secret type IDs
		for key, val := range secretIDs {
			resp, err := client.R().
				SetBody(fmt.Sprintf(`{"id": "%s"}`, key)).
				Post(fmt.Sprintf("https://vault.%s/v1/get", pangeaDomain))

			if err != nil {
				log.Fatal("Error fetching secret ", val)
			}

			var response VaultSecretResponse
			err = json.Unmarshal(resp.Body(), &response)

			if err != nil {
				log.Fatalf("Error while fetching secrets. %s", err)
			}

			remoteEnv = append(remoteEnv, fmt.Sprintf("%s=%s", val, response.Result.CurrentVersion.Secret))
		}
	} else {
		log.Fatal("Folder not found. Please use `pangea select` to choose the workspace you would like to use secrets from.")
	}

	return remoteEnv
}

func exec_subprocess(baseCommand string, args []string) error {
	cmd := exec.Command(baseCommand, args...)

	remoteEnv := Get_env()

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
	rootCmd.AddCommand(runCmd)
}
