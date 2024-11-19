package updates

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/pangeacyber/pangea-cli/v2/plugins"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var PluginCheckUpdate = plugins.NewPlugin(checkUpdateCmd, []string{"check-update"})
var PluginUpdate = plugins.NewPlugin(updateCmd, []string{"update"})

// Avoid to print multiple times
var availablePrinted bool

func CheckAvailableVersion(skipCache bool) (*cli.Release, bool, error) {
	found := false
	var latestSelfupdate *selfupdate.Release
	var err error
	var latest *cli.Release

	if !skipCache {
		// First check release on cache
		latest, _ = cli.CacheGetRelease()
	}

	// if no release data on cache, request it.
	if latest == nil {
		latestSelfupdate, found, err = selfupdate.DetectLatest("pangeacyber/pangea-cli")
		if err != nil {
			return nil, false, fmt.Errorf("error occurred while detecting version: %v", err)
		}
	}

	// If requested release and it was found. Save it to cache.
	if found {
		latest = cli.NewRelease(latestSelfupdate)
		_ = cli.CacheSetRelease(latest)
	}

	if latest == nil {
		return nil, false, nil
	}

	currentVersion := semver.MustParse(cli.Version[1:])
	latestVersion := semver.MustParse(latest.Version)
	if latestVersion.LTE(currentVersion) {
		return nil, false, nil
	}

	if !availablePrinted {
		fmt.Fprintln(os.Stderr, "New version available:", latest.Version)
		availablePrinted = true
	}

	return latest, true, nil
}

func printAvailableVersionInfo(skipCache bool) error {
	latest, available, err := CheckAvailableVersion(skipCache)
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

func printReleaseInfo(latest *cli.Release) {
	if latest == nil {
		return
	}
	fmt.Println("Release notes:\n", latest.ReleaseNotes)
}

func update() error {
	latest, available, err := CheckAvailableVersion(true)
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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Pangea CLI to the latest version",
	Long:  `Update Pangea CLI to the latest version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return update()
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check-update",
	Short: "Check if there is a new Pangea CLI version",
	Long:  `Check if there is a new Pangea CLI version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skipCache, err := cmd.Flags().GetBool("skip-cache")
		if err != nil {
			return err
		}

		return printAvailableVersionInfo(skipCache)
	},
}

func init() {
	checkUpdateCmd.Flags().BoolP("skip-cache", "s", false, "Set to true to skip cache data and force request to Github")
}
