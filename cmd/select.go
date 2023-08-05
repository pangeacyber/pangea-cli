/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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
	Short: "Select the project you want to get secrets from on Pangea",
	Long: `This command selects the project you want to link your current directory to a remote directory on Pangea vault.
	You need to do this before you use "pangea run" to specify which directory you want to fetch secrets from.`,
	Run: func(cmd *cobra.Command, args []string) {
		remoteFolderName, err := cmd.Flags().GetString("folder-name")
		if err != nil {
			log.Fatal(err)
		}

		SelectProject(remoteFolderName)
	},
}

func SelectProject(remoteFolderName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	cachePath := utils.GetCachePath()
	config := utils.LoadCacheData(cachePath)

	if remoteFolderName == "" {
		projectName := promptUser("Enter the name of your project: ")

		projectEnv := SelectProjectEnvironment()

		config.Paths[currentDir] = utils.ProjectData{
			Remote: "/secrets/" + projectName + "/" + projectEnv,
		}

		saveCacheData(cachePath, config)

		fmt.Printf("Project '%s' selected and cached.\n", projectName)
	} else {
		config.Paths[currentDir] = utils.ProjectData{
			Remote: remoteFolderName,
		}

		saveCacheData(cachePath, config)

		fmt.Printf("Project '%s' selected and cached.\n", remoteFolderName)
	}
}

func promptUser(promptMessage string) string {
	fmt.Print(promptMessage)
	var input string
	fmt.Scanln(&input)
	return input
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func SelectProjectEnvironment() string {
	return promptUser("Which environment would you like to use (dev / stg / prod): ")
}

func saveCacheData(cachePath string, config utils.CacheData) {
	viper.SetConfigFile(cachePath)
	viper.SetConfigType("json")

	viper.ReadInConfig()
	// if err != nil {
	// 	fmt.Printf("Error reading cache file: %v\n", err)
	// }

	viper.Set("paths", config.Paths)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing cache file: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(selectCmd)

	selectCmd.Flags().StringP("folder-name", "f", "", "folder name on Pangea vault (Example - /secrets/<project_name>/dev)")
}
