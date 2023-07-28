/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
)

var command []string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run your application with secrets on Pangea",
	Long: `Run your applications with secrets loaded as environment variables into your application on Pangea.
	
	For example:
		pangea run -c npm run dev
			- will start your node server with secrets loaded in from Pangea`,
	Run: func(cmd *cobra.Command, args []string) {
		baseCommand, err := cmd.Flags().GetStringArray("command")
		if err != nil {
			log.Fatal("No specified command")
		}

		err = exec_subprocess(baseCommand, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func get_env() {
	cachePath := GetCachePath()
	config := LoadCacheData(cachePath)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error reading current directory path")
	}

	if _, isPathExists := config.Paths[currentDir]; isPathExists {
		client := utils.CreateVaultAPIClient()

		_, err := client.R().
			// SetBody(fmt.Sprintf(`{"name":"%s", "folder":"%s"}`, env, fmt.Sprintf("secrets/%s/", projectName))).
			Post("https://vault.aws.us.pangea.cloud/v1/folder/create")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func exec_subprocess(baseCommand []string, args []string) error {
	get_env()
	cmd := exec.Command(baseCommand[0], args...)

	env := make([]string, len(os.Environ()))
	copy(env, os.Environ())
	env = append(env, "TYPE=PROD")

	cmd.Env = env

	// Set up pipes for stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the subprocess
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting subprocess:", err)
		return err
	}

	// Wait for the subprocess to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for subprocess:", err)
		return err
	}

	// Subprocess has completed

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringArrayVarP(&command, "command", "c", []string{}, "Command to execute")
}
