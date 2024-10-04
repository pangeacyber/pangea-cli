// Copyright 2023 Pangea Cyber Corporation
// Author: Pangea Cyber Corporation

package authentication

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/pangeacyber/pangea-cli-internal/cli"
	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/spf13/cobra"
)

var PluginLogin = plugins.NewPlugin(loginCmd, []string{"login"})

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

func loginPrompts(noBrowserStatus bool, profile, service string) {
	if !noBrowserStatus {
		fmt.Println("When you hit enter, we will redirect you to the Pangea Vault page where you will need to copy the Default Pangea Token and paste it in the next prompt.")
		_, _ = fmt.Scanln()

		err := openBrowser("https://console.pangea.cloud/service/vault")
		if err != nil {
			// TODO: Handle error if browser doesn't work
			log.Fatal(err)
		}
	} else {
		fmt.Println("Visit https://console.pangea.cloud/service/vault to grab your Default Pangea Token and paste it below.")
	}

	fmt.Print("Enter Pangea Token (without any quotes): ")
	token := cli.ReadStdin()
	if token == "" {
		log.Fatal("Invalid empty token")
	}

	fmt.Print("Enter Pangea Domain (without any quotes): ")
	domain := cli.ReadStdin()

	if profile == "" {
		fmt.Print("Enter profile name to associate this token and domain. Press enter to use 'default' profile: ")
		profile = cli.ReadStdin()
	}

	if service == "" {
		fmt.Print("Enter service name to associate this token and domain. Press enter to use 'default' service: ")
		service = cli.ReadStdin()
	}

	err := cli.SaveToken(profile, service, token)
	if err != nil {
		log.Fatal(err)
	}

	err = cli.SaveDomain(profile, service, domain)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Save config success.")
}

var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Connect Pangea CLI to Pangea Vault",
	Long:    `Login to Pangea CLI and connect it to Pangae's Vault.`,
	GroupID: "tools",
	Run: func(cmd *cobra.Command, args []string) {
		noBrowser, _ := cmd.Flags().GetBool("no-browser")
		service, _ := cmd.Flags().GetString("service")
		profile, _ := cmd.Flags().GetString("profile")
		loginPrompts(noBrowser, profile, service)
	},
}

func init() {
	// loginCmd represents the login command
	loginCmd.Flags().Bool("no-browser", false, "Login without the browser")
	loginCmd.Flags().StringP("profile", "p", "", "Associate token to 'profile'")
	loginCmd.Flags().StringP("service", "s", "", "Associate token to 'service'")
}
