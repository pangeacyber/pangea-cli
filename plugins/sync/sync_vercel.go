package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pangeacyber/pangea-cli-internal/cli"
	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/pangeacyber/pangea-cli-internal/plugins/vault"
	"github.com/spf13/cobra"
)

const vercelVersion = "v9"

var (
	vercelToken     string
	vercelProjectId string
	vercelTarget    string
	vercelGitBranch string
)

var PluginSyncVercel = plugins.NewPlugin(vercelCmd, []string{"sync", "vercel"})
var logger = cli.GetLogger()

var vercelCmd = &cobra.Command{
	Use:   "vercel",
	Short: "Sync secrets from Vault workspace to Vercel.",
	Long:  "Sync secrets from Vault workspace to Vercel.",
	Run: func(cmd *cobra.Command, args []string) {
		if vercelToken == "" {
			vercelToken = os.Getenv("VERCEL_TOKEN")
		}

		if vercelProjectId == "" {
			vercelProjectId = os.Getenv("VERCEL_PROJECT_ID")
		}

		if vercelToken == "" || vercelProjectId == "" {
			logger.Fatal("Vercel token and project ID must be provided either as flags or environment variables.")
		}

		envs := vault.GetWorkspaceSecrets(vault.GetWorkspaceFromSettings())
		if err := pushEnvToVercel(envs); err != nil {
			logger.Fatalf("Error while syncing secrets from Vault workspace to vercel: %s\n", err.Error())
		}
	},
}

func init() {
	syncCmd.AddCommand(vercelCmd)
	vercelCmd.Flags().StringVarP(&vercelToken, "token", "t", "", "Vercel API token")
	vercelCmd.Flags().StringVarP(&vercelProjectId, "project", "p", "", "Vercel project ID")
	vercelCmd.Flags().StringVarP(&vercelTarget, "target", "x", "development", "Comma separated list of vercel environments to push to, defaults to 'development'")
	vercelCmd.Flags().StringVarP(&vercelGitBranch, "branch", "b", "", "Which git branch to allow access to this variable, target must be set to 'preview'.")
}

func pushEnvToVercel(envs []string) error {
	envURL := fmt.Sprintf("https://api.vercel.com/%s/projects/%s/env", vercelVersion, vercelProjectId)
	environment := make(map[string]string, len(envs))
	for _, value := range envs {
		name, val, found := strings.Cut(value, "=")
		if !found {
			logger.Fatalf("Variable %s is invalid", value)
		}
		environment[name] = val
	}

	if len(environment) == 0 {
		logger.Println("No secrets to push")
		return nil
	}

	client := &http.Client{}
	tokenHeaderValue := fmt.Sprintf("Bearer %s", vercelToken)

	// Fetch project secrets
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?source=%s", envURL, "pangea"), nil)
	req.Header.Set("Authorization", tokenHeaderValue)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatalf("Failed to fetch variables from vercel due to: '%s'\n", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var vercelEnvs VercelEnvResponse
	if err := dec.Decode(&vercelEnvs); err != nil {
		logger.Fatalln("Failed to parse response from vercel due to: ", err)
	}

	// Update secrets
	var targets []string
	if vercelTarget != "" {
		targets = strings.Split(vercelTarget, ",")
	}

	if vercelGitBranch != "" {
		if targets == nil {
			targets = []string{"preview"}
		}
		if len(targets) != 1 || targets[0] != "preview" {
			logger.Fatalln("If vercel branch is provided, the target must be 'preview'")
		}
	}

	if targets == nil {
		targets = []string{"development"}
	}

	for key, value := range environment {

		envVar := map[string]any{
			"key":   key,
			"value": value,
			"type":  "encrypted",
		}
		if targets != nil {
			envVar["target"] = targets
		}
		if vercelGitBranch != "" {
			envVar["gitBranch"] = vercelGitBranch
		}
		body, _ := json.Marshal(envVar)
		var req *http.Request
		if id, ok := vercelEnvs.KeyExists(key); ok {
			// Update variable
			req, err = http.NewRequest("PATCH", fmt.Sprintf("%s/%s", envURL, id), bytes.NewBuffer(body))
		} else {
			// Create new variable
			req, err = http.NewRequest("POST", envURL, bytes.NewBuffer(body))
		}

		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		req.Header.Set("Authorization", tokenHeaderValue)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error making request: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to push env var %s, received status: %d: %s", key, resp.StatusCode, string(b))
		}
		logger.Printf("Successfully pushed %s\n", key)
	}
	return nil
}

type VercelEnvResponse struct {
	Envs []struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	} `json:"envs"`
}

func (v *VercelEnvResponse) KeyExists(key string) (string, bool) {
	for _, env := range v.Envs {
		if env.Key == key {
			return env.ID, true
		}
	}
	return "", false
}
