/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CacheData struct {
	Paths map[string]ProjectData `json:"paths"`
}

type ProjectData struct {
	Remote string `json:"remote"`
}

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

		cachePath := getCachePath()
		config := loadCacheData(cachePath)

		config.Paths[currentDir] = ProjectData{
			Remote: "/secrets/" + projectName,
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

func getCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".pangea", "cache_paths.json")
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func loadCacheData(cachePath string) CacheData {
	viper.SetConfigFile(cachePath)
	viper.SetConfigType("json")
	viper.SetDefault("paths", make(map[string]ProjectData))

	err := viper.ReadInConfig()
	if err != nil {
		// fmt.Printf("Error reading cache file: %v\n", err)
	}

	var config CacheData
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Error parsing cache file: %v\n", err)
	}

	return config
}

func saveCacheData(cachePath string, config CacheData) {
	viper.SetConfigFile(cachePath)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading cache file: %v\n", err)
	}

	viper.Set("paths", config.Paths)

	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing cache file: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(selectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
