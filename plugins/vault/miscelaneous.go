package vault

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	sv "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
)

var logger = cli.GetLogger()

// Get current directory
func GetWD() string {
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal("Error reading current directory path")
	}
	return strings.ToLower(wd)
}

func CreateVaultService() (sv.Client, error) {
	token, domain, err := cli.GetTokenAndDomain("vault")
	if err != nil {
		return nil, err
	}
	config := cli.GetDefaultPangeaConfig()
	config.Token = token
	config.Domain = domain

	return sv.New(&config), nil
}

func GetWorkspaceFromSettings() string {
	defaultWorkspace := os.Getenv("PANGEA_DEFAULT_FOLDER")
	workspace := ""

	wd := GetWD()

	config, err := cli.LoadCacheData()
	if err != nil {
		return ""
	}

	w, ok := config.Paths[wd]
	if ok {
		workspace = w.Remote
	}

	if workspace == "" {
		workspace = defaultWorkspace
	}

	return workspace
}

func GetWorkspaceSecrets(workspace string) []string {
	var remoteEnv []string

	if workspace != "" {
		logger.Printf("Fetching secrets from: %s\n", workspace)
		client, err := CreateVaultService()
		if err != nil {
			return []string{}
		}

		resp, err := client.List(
			context.Background(),
			&sv.ListRequest{
				Filter: map[string]string{
					"folder": workspace,
				},
			},
		)
		if err != nil {
			logger.Fatal("Error fetching secrets from Pangea")
		}

		if resp.Status != nil && *resp.Status == "Unauthorized" {
			logger.Fatal("Unauthorized! Please run `pangea login` to get a new token.")
		}

		// Create a list of secret type IDs
		secretIDs := make(map[string]interface{})
		for _, item := range resp.Result.Items {
			if item.Type == "secret" {
				secretIDs[item.ID] = item.Name
			}
		}

		// Print the list of secret type IDs
		for id, val := range secretIDs {
			resp, err := client.Get(
				context.Background(),
				&sv.GetRequest{
					ID: id,
				},
			)
			if err != nil {
				logger.Fatal("Error fetching secret ", val)
			}
			if resp.Result.CurrentVersion.Secret != nil {
				remoteEnv = append(remoteEnv, fmt.Sprintf("%s=%s", val, *resp.Result.CurrentVersion.Secret))
			}
		}
	} else {
		logger.Fatal("Folder not found. Please use `pangea vault workspace select` to choose the workspace you would like to use secrets from")
	}

	return remoteEnv
}
