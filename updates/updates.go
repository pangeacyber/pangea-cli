package updates

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/pangeacyber/pangea-cli/cmd"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// Avoid to print multiple times
var availablePrinted bool

func CheckAvailableVersion() (*selfupdate.Release, bool, error) {
	currentVersion := semver.MustParse(cmd.Version[1:])

	latest, found, err := selfupdate.DetectLatest("pangeacyber/pangea-cli")
	if err != nil {
		return nil, false, fmt.Errorf("error occurred while detecting version: %v", err)
	}

	if !found || latest.Version.LTE(currentVersion) {
		return nil, false, nil
	}

	if !availablePrinted {
		fmt.Fprintln(os.Stderr, "New version available:", latest.Version)
		availablePrinted = true
	}

	return latest, true, nil
}

func PrintAvailableVersionInfo() error {
	latest, available, err := CheckAvailableVersion()
	if err != nil {
		return err
	}

	if !available {
		fmt.Println("Current version is the latest")
		return nil
	}

	printReleaseInfo(latest)
	return nil
}

func printReleaseInfo(latest *selfupdate.Release) {
	if latest == nil {
		return
	}
	fmt.Println("Release notes:\n", latest.ReleaseNotes)
}

func Update() error {
	latest, available, err := CheckAvailableVersion()
	if err != nil {
		return err
	}

	if !available {
		fmt.Println("Current version is the latest")
		return nil
	}

	printReleaseInfo(latest)

	// Get the absolute path of the executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	fmt.Printf("Updating executable: %s\n", exePath)

	if err := selfupdate.UpdateTo(latest.AssetURL, exePath); err != nil {
		return fmt.Errorf("error occurred while updating binary: %v", err)
	}
	return nil
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Pangea CLI to the latest version",
	Long:  `Update Pangea CLI to the latest version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Update()
	},
}

var CheckUpdateCmd = &cobra.Command{
	Use:   "check-update",
	Short: "Check if there is a new Pangea CLI version",
	Long:  `Check if there is a new Pangea CLI version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return PrintAvailableVersionInfo()
	},
}

func init() {
	cmd.RootCmd.AddCommand(UpdateCmd)
	cmd.RootCmd.AddCommand(CheckUpdateCmd)
}
