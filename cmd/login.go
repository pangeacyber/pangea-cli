/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/pangeacyber/pangea-cli/utils"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Connect CLI to Pangea Vault",
	Long:  `Login to CLI and connect it to Pangae's Vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		noBrowser, _ := cmd.Flags().GetBool("no-browser")
		loginPrompts(noBrowser)
	},
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func loginPrompts(noBrowserStatus bool) {
	if !noBrowserStatus {
		fmt.Println("When you hit enter, we will redirect you to the Pangea Vault page where you will need to copy the Default Pangea Token and paste it in the next prompt. Ready 🏎️?")
		fmt.Scanln()

		err := openBrowser("https://console.pangea.cloud/service/vault")
		if err != nil {
			// TODO: Handle error if browser doesn't work
			log.Fatal(err)
		}
	} else {
		fmt.Println("Visit https://console.pangea.cloud/service/vault to grab your Default Pangea Token and paste it below.")
	}

	fmt.Print("Enter Pangea Token: ")
	// TODO: Check if Pangea token is valid
	token := readInput()

	fmt.Print("Enter Pangea Domain: ")
	domain := readInput()

	err := utils.WriteTokenToFile(token, domain)
	if err != nil {
		log.Fatal(err)
	}

}

// readInput reads user input from the command line
func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().Bool("no-browser", false, "Login without the browser")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
