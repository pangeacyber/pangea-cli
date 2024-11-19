package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	pangea "github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/pangeacyber/pangea-cli/v2/builder"
	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/pangeacyber/pangea-cli/v2/plugins/authentication"
	loader "github.com/pangeacyber/pangea-cli/v2/plugins/plugins_loader"
	"github.com/pangeacyber/pangea-cli/v2/plugins/updates"
)

const (
	OpenAPIPath = "/v1/openapi.json"
)

var Services = []string{
	"vault",
	"embargo",
	"redact",
	"file-intel",
}

var (
	rootCmd = &cobra.Command{
		Use:   "pangea",
		Short: "Pangea CLI " + cli.Version,
	}

	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print Pangea CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cli.Version)
		},
	}

	// versionCmd represents the version command
	cleanCacheCmd = &cobra.Command{
		Use:   "clean",
		Short: "Remove cached json schema files",
		RunE: func(cmd *cobra.Command, args []string) error {
			folder, err := cli.GetCacheFolder()
			if err != nil {
				return err
			}
			err = os.RemoveAll(folder)
			if err != nil {
				fmt.Printf("Failed to delete %s folder. Error: %v\n", folder, err)
			}
			fmt.Printf("Deleted folder: %s.\n", folder)
			return nil
		},
	}

	adminCmd = &cobra.Command{
		Use:     "admin",
		Short:   "List of 'admin' commands.",
		GroupID: "tools",
	}

	utilsCmd = &cobra.Command{
		Use:     "utils",
		Short:   "List of 'utils' commands.",
		GroupID: "tools",
	}

	loadedServices = map[string]bool{"": true}
)

func main() {
	_, _, err := updates.CheckAvailableVersion(false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rootCmd.AddGroup(
		&cobra.Group{
			Title: "Services",
			ID:    "services",
		},
		&cobra.Group{
			Title: "Tools",
			ID:    "tools",
		},
	)

	b := builder.NewBuilder(rootCmd)
	err = b.AddCommand([]string{"version"}, versionCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to 'version' command. Error: %v", err)
	}

	err = b.AddCommand([]string{"admin"}, adminCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to 'admin' command. Error: %v", err)
	}

	err = b.AddCommand([]string{"utils"}, utilsCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to 'utils' command. Error: %v", err)
	}

	err = b.AddCommand([]string{"admin", "cache", "clean"}, cleanCacheCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to 'admin cache clean' command. Error: %v", err)
	}

	processServices(b)
	processPlugins(b)

	// Always load PluginLogin so user can authenticate
	err = b.AddPlugin(authentication.PluginLogin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add PluginLogin. Error: %v.\n", err)
	}

	err = b.Build(rootCmd)
	if err != nil {
		log.Fatal(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing root command. %v", err)
	}
}

func processServices(b *builder.Builder) {
	errorPrinted := false
	for _, svc := range Services {
		token, domain, err := cli.GetTokenAndDomain(svc)
		if err != nil {
			if errorPrinted {
				continue
			}
			fmt.Fprintf(os.Stderr, "\nError loading Pangea Token from config file. Error:[%v]\nCreate or update a profile running 'pangea admin profile' command\n\n", err)
			errorPrinted = true
			continue
		}

		config := cli.GetDefaultPangeaConfig()
		config.Domain = domain
		config.Token = token

		svcCmd := &cobra.Command{
			Use:     svc,
			Short:   fmt.Sprintf("Pangea %s service", cases.Title(language.English).String(svc)),
			GroupID: "services",
		}

		svcCmd.PersistentFlags().String("cli-profile", "", "Run API call with a particular CLI profile token and domain. Setup profile with 'pangea admin profile' commands.")

		err = processServiceJsonSchema(b, svc, &config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to process %s service json schema. Error: %v.\n", svc, err)
			continue
		}

		err = b.AddCommand([]string{svc}, svcCmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add service[%s]. Error: %v", svc, err)
			continue
		}

		loadedServices[svc] = true
	}
}

func processPlugins(b *builder.Builder) {
	for _, p := range loader.LoadPlugins() {
		as := p.GetServiceAssociated()
		if as != "" && !loadedServices[as] {
			continue
		}

		err := p.Init(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Plugin %s failed to init: %v.\n", reflect.TypeOf(p), err)
			continue
		}

		err = b.AddPlugin(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add plugin %s. Error: %v.\n", reflect.TypeOf(p), err)
			continue
		}
	}
}

func processServiceJsonSchema(b *builder.Builder, svc string, config *pangea.Config) error {
	client := pangea.NewClient(svc, config)

	oapiURL, err := client.GetURL(OpenAPIPath)
	if err != nil {
		return err
	}

	oapi, err := cli.LoadURL(oapiURL)
	if err != nil {
		return err
	}

	if oapi.Status != nil && *oapi.Status == "Unauthorized" {
		fmt.Fprintln(os.Stderr, "Unauthorized token. Execute: 'pangea login' to update it.")
		err = cli.RemoveCachedFileFromURL(oapiURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove cached file from URL. Error: %v", err)
		}

		return cli.ErrUnauthorized
	}

	// add all commands (keep version version)
	for path, pathDef := range oapi.Paths {
		pathOri := path
		version := ""
		ndx := strings.Index(path[1:], "/")
		if ndx >= 0 {
			version = path[1 : ndx+1]
		}
		b.AddPangeaCommand(config, svc, path[ndx+1:], pathOri, pathDef.Post, version)
	}
	return nil
}
