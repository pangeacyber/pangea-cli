// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select the workspace you want to get secrets from on Pangea Vault",
	Long: `This command selects the workspace you want to link your current directory to a remote directory on Pangea Vault.
	You need to do this before you use "pangea run" to specify which directory you want to fetch secrets from.`,
	Run: func(cmd *cobra.Command, args []string) {
		remoteFolderName, err := cmd.Flags().GetString("folder-name")
		if err != nil {
			log.Fatal(err)
		}

		SelectWorkspace(remoteFolderName)
	},
}

func SelectWorkspace(remoteFolderName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	cachePath := utils.GetCachePath()
	config := utils.LoadCacheData(cachePath)

	if remoteFolderName == "" {
		workspaceName := promptUser("Enter the name of your workspace: ")

		// Removing env selection for v1
		// workspaceEnv := SelectWorkspaceEnvironment()

		config.Paths[currentDir] = utils.WorkspaceData{
			Remote: "/secrets/" + workspaceName,
		}

		saveCacheData(cachePath, config)

		fmt.Printf("Workspace '%s' selected and cached.\n", workspaceName)
	} else {
		config.Paths[currentDir] = utils.WorkspaceData{
			Remote: remoteFolderName,
		}

		saveCacheData(cachePath, config)

		fmt.Printf("Workspace '%s' selected and cached.\n", remoteFolderName)
	}
}

func promptUser(promptMessage string) string {
	fmt.Print(promptMessage)
	var input string
	fmt.Scanln(&input)
	return input
}

func SelectWorkspaceEnvironment() string {
	return promptUser("Which environment would you like to use (dev / stg / prod): ")
}

func saveCacheData(cachePath string, config utils.CacheData) {
	viper.SetConfigFile(cachePath)
	viper.SetConfigType("json")

	viper.ReadInConfig() //nolint:errcheck
	// if err != nil {
	// 	fmt.Printf("Error reading cache file: %v\n", err)
	// }

	viper.Set("paths", config.Paths) //nolint:errcheck

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing cache file: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(selectCmd)

	selectCmd.Flags().StringP("folder-name", "f", "", "folder name on Pangea vault (Example - /secrets/<workspace_name>/dev)")
}
