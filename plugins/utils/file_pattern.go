package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pangeacyber/pangea-cli-internal/builder"
	"github.com/pangeacyber/pangea-cli-internal/cli"
	"github.com/pangeacyber/pangea-cli-internal/plugins"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const cmdUse = "file-pattern"

type FileProcess func(cmd *cobra.Command, filesNames []string) (map[string]string, error)
type ResponseProcess func(filesToValues map[string]string, response map[string]any) error

type pluginWithFilePattern struct {
	pluginCmd         *cobra.Command
	cmdPath           []string
	pluginPath        []string
	flagName          string
	callback          FileProcess
	responseProcess   ResponseProcess
	fileToValue       map[string]string
	serviceAssociated string
}

func NewPluginWithFilePattern(cmdPath []string, flagName string, callback FileProcess, responseProcess ResponseProcess, serviceAssociated string) plugins.Plugin {
	pluginPath := append(cmdPath, cmdUse)
	p := &pluginWithFilePattern{
		cmdPath:           cmdPath,
		pluginPath:        pluginPath,
		flagName:          flagName,
		callback:          callback,
		fileToValue:       map[string]string{},
		responseProcess:   responseProcess,
		serviceAssociated: serviceAssociated,
	}
	return p
}

func (p *pluginWithFilePattern) GetCommandPath() []string {
	return p.pluginPath
}

func (p *pluginWithFilePattern) GetCommand() *cobra.Command {
	return p.pluginCmd
}

func (p *pluginWithFilePattern) Init(b plugins.CommandProvider) error {
	if p.cmdPath == nil {
		return errors.New("`cmdPath` should not be nil")
	}
	if p.callback == nil {
		return errors.New("`callback` should not be nil")
	}

	cmd := b.GetCommand(p.cmdPath)
	if cmd == nil {
		return fmt.Errorf("not available command [%s]", plugins.CommandPathToString(p.cmdPath))
	}

	flag := cmd.Flag(p.flagName)
	if flag == nil {
		return fmt.Errorf("`%s` does not exist in base command", p.flagName)
	}

	flagValue, ok := flag.Value.(*builder.FlagArray)
	if !ok {
		return fmt.Errorf("`%s` flag should be of type `%s` and is type `%s`", p.flagName, reflect.TypeOf(builder.FlagArray{}), reflect.TypeOf(flag))
	}

	// mark flags as optional
	flags := cmd.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		_ = flags.SetAnnotation(f.Name, cobra.BashCompOneRequiredFlag, []string{"false"})
	})

	cmd.Use = cmdUse
	cmd.Short = cmd.Short + " Using file pattern"
	cmd.Long = cmd.Long + " Using file pattern"
	cmd.Args = cobra.ExactArgs(1)

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Println("Usage:")
		fmt.Printf("    pangea %s \"*.txt,*.json,pangea*\" [flags]\n", plugins.CommandPathToString(p.GetCommandPath()))
		fmt.Println("\nFlags:")
		fmt.Println(cmd.Flags().FlagUsages())
		return nil
	})

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		patterns, err := cli.ReadAsCSV(args[0])
		if err != nil {
			return err
		}
		files := []string{}
		patterns = replaceUserFolder(patterns)

		for _, pattern := range patterns {
			newFiles, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error looking for file pattern: `%s`. %v", pattern, err)
				return err
			}
			files = append(files, newFiles...)
		}

		if len(files) == 0 {
			return errors.New("no files match given pattern")
		}

		p.fileToValue, err = p.callback(cmd, files)
		if err != nil {
			return err
		}

		values := make([]string, 0, len(p.fileToValue))
		for _, value := range p.fileToValue {
			values = append(values, value)
		}

		_ = flagValue.Replace(values)
		flag.Changed = true
		return nil
	}

	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		if p.responseProcess == nil {
			return nil
		}
		if cmd.Annotations == nil || cmd.Annotations["pangea_response"] == "" {
			return p.responseProcess(p.fileToValue, nil)
		}
		respData := map[string]any{}
		err := json.Unmarshal([]byte(cmd.Annotations["pangea_response"]), &respData)
		if err != nil {
			return err
		}

		return p.responseProcess(p.fileToValue, respData)
	}

	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	p.pluginCmd = cmd
	return nil
}

func (p *pluginWithFilePattern) GetServiceAssociated() string {
	return p.serviceAssociated
}

func replaceUserFolder(patterns []string) []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return patterns
	}

	newPatterns := []string{}

	for _, p := range patterns {
		if strings.HasPrefix(p, "~/") {
			p = filepath.Join(homeDir, p[2:])
		}
		newPatterns = append(newPatterns, p)
	}
	return newPatterns
}
