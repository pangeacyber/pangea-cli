/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := promptUser("Enter the name of your project: ")
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}

		cachePath := utils.GetCachePath()
		config := utils.LoadCacheData(cachePath)

		projectEnv := SelectProjectEnvironment()

		config.Paths[currentDir] = utils.ProjectData{
			Remote: "/secrets/" + projectName + "/" + projectEnv,
		}

		saveCacheData(cachePath, config)

		fmt.Printf("Project '%s' selected and cached.\n", projectName)
	},
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
}
