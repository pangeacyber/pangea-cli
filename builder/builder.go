package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	pe "github.com/huantt/plaintext-extractor"
	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/pangeacyber/pangea-cli/v2/plugins"
	pangea "github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Builder struct {
	cmdMap *CommandsMap
}

type CommandsMap struct {
	Command     *cobra.Command
	SubCommands map[string]*CommandsMap
}

func NewBuilder(cmd *cobra.Command) *Builder {
	return &Builder{
		cmdMap: &CommandsMap{
			Command: cmd,
		},
	}
}

func (b *Builder) Cmd() *cobra.Command {
	return b.cmdMap.Command
}

// AddPangeaCommand creates a new command from an openapi path spec.
func (b *Builder) AddPangeaCommand(config *pangea.Config, svc, pathCmd, pathAPI string, post cli.PathPost, version string) {
	if post.XPangeaUISchema != nil && post.XPangeaUISchema.IsConfiguration != nil && *post.XPangeaUISchema.IsConfiguration {
		return
	}

	groupID := ""
	if post.Tags != nil && len(post.Tags) == 1 {
		groupID = post.Tags[0]
	}

	short := post.Summary
	if len(short) < 30 { // If it's too short, use description
		short = strings.Split(post.Description, "\n")[0]
	}

	cmd := &cobra.Command{
		Use:     pathCmd,
		Short:   cleanFormat(short),
		Long:    cleanFormat(post.Description),
		RunE:    b.executePangeaRequest(config, svc, pathAPI),
		GroupID: groupID,
		Annotations: map[string]string{
			"version": version,
		},
	}

	addSchema(cmd, post.RequestBody.Content[cli.ApplicationJSON].Schema, true)
	err := b.AddCommand([]string{svc, version, pathCmd}, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add service[%s] command path[%s]. Error: %v", svc, pathCmd, err)
	}
}

func (b *Builder) AddCommand(path []string, cmd *cobra.Command) error {
	if len(path) >= 1 {
		name := path[len(path)-1]
		if name != cmd.Use {
			return fmt.Errorf("last path name [%s] and command `Use` [%s] does not match", name, cmd.Use)
		}
	}

	return b.cmdMap.addSubCommand(path, cmd)
}

func (b *Builder) Build(rootCmd *cobra.Command) error {
	if b.cmdMap.Command == nil {
		return errors.New("base command should not be nil")
	}
	b.build(rootCmd, b.cmdMap.Command.Name(), b.cmdMap)
	return nil
}

func (b *Builder) AddPlugin(p plugins.Plugin) error {
	return b.AddCommand(p.GetCommandPath(), p.GetCommand())
}

func (b *Builder) GetCommand(path []string) *cobra.Command {
	// Return a copy of the command in that path if exists, any other case return nil
	return b.cmdMap.getCommand(path)
}

// propEnumValues returns a list of possible values for a property, if it has an enum or const field.
// It returns an empty list if the property has no enum or const field.
func propEnumValues(prop cli.Property) []string {
	var vals []string
	for _, val := range prop.Enum {
		if vs, ok := val.(string); ok {
			vals = append(vals, vs)
		}
	}
	if prop.Const != nil {
		if c, ok := prop.Const.(string); ok {
			vals = append(vals, c)
		}
	}
	return vals
}

func isConstStringEnum(prop cli.Property) bool {
	switch prop.Const.(type) {
	case string:
		return true
	default:
		return false
	}
}

func mergeDescriptions(oldDesc, newDesc string) string {
	updateDesc := ""
	oldDesc = strings.TrimSpace(oldDesc)
	if oldDesc != "" && !strings.HasSuffix(oldDesc, ".") {
		updateDesc = oldDesc + ". "
	} else if oldDesc != "" && strings.HasSuffix(oldDesc, ".") {
		updateDesc = oldDesc + " "
	}
	return updateDesc + newDesc
}

// addParameters add a list of properties as flags to a command
// WIP
func addParameters(cmd *cobra.Command, props cli.Properties, required []string) {
	for name, prop := range props {
		propType := string(prop.Type)

		if strings.Contains(propType, `"string"`) || isConstStringEnum(prop) {
			var def string
			def, _ = prop.Default.(string)
			vals := propEnumValues(prop)
			if len(vals) > 0 {
				var fe *FlagEnum
				if flag := cmd.Flag(name); flag != nil {
					fe = flag.Value.(*FlagEnum)
					fe.AddValues(vals)
					flag.Usage = mergeDescriptions(prop.Description, fe.Description())
				} else {
					fe = NewFlagEnum(name, vals)
					descr := mergeDescriptions(prop.Description, fe.Description())
					cmd.Flags().Var(fe, name, descr)
				}

				// add the flag values to the autocomplete scripts
				cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
					return fe.GetValues(), cobra.ShellCompDirectiveDefault
				})
				continue
			}

			if cmd.Flags().Lookup(name) != nil {
				// TODO
				continue
			}

			cmd.Flags().String(name, def, prop.Description)
			continue
		}

		if cmd.Flags().Lookup(name) != nil {
			// TODO: What we should do?
			continue
		}

		if strings.Contains(propType, `"integer"`) {
			var f FlagInteger
			cmd.Flags().Var(&f, name, prop.Description)
			continue
		}

		if strings.Contains(propType, `"boolean"`) {
			var f FlagBool
			cmd.Flags().Var(&f, name, prop.Description)
			continue
		}

		if strings.Contains(propType, `"object"`) {
			var f FlagMap
			help := fmt.Sprintf("CLI use: '--%s key1:value1,key2:value2'.", name)
			cmd.Flags().Var(&f, name, cleanFormat(mergeDescriptions(prop.Description, help)))
			continue
		}

		if strings.Contains(propType, `"array"`) {
			var f FlagArray
			help := fmt.Sprintf("CLI use: '--%s value1,value2'.", name)
			cmd.Flags().Var(&f, name, cleanFormat(mergeDescriptions(prop.Description, help)))
			continue
		}

		// By default add flag as Any. Could be string or an object. It apply to some `redact fields`
		var f FlagAny
		cmd.Flags().Var(&f, name, cleanFormat(prop.Description))
	}

	for _, flag := range required {
		cmd.MarkFlagRequired(flag) //nolint:errcheck
		// err := cmd.MarkFlagRequired(flag)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to mark flag as required: %v\n", err)
		// }
	}
}

func cleanMD(in string) string {
	emd := pe.NewMarkdownExtractor()
	r, err := emd.PlainText(in)
	if err != nil {
		return ""
	}
	return *r
}

func cleanHTML(in string) string {
	ehtml := pe.NewHtmlExtractor()
	r, err := ehtml.PlainText(in)
	if err != nil {
		return ""
	}
	return *r
}

func cleanFormat(in string) string {
	return cleanHTML(cleanMD(in))
}

// addSchema populates a command with the attributes of a schema.
// It follows oneOf and anyOf references.
// WIP
func addSchema(cmd *cobra.Command, schema *cli.Schema, wr bool) {
	required := []string{}
	if wr {
		required = schema.Required
	}

	addParameters(cmd, schema.Properties, required)
	for _, sch := range schema.OneOf {
		addSchema(cmd, &sch, false)
	}
	for _, sch := range schema.AnyOf {
		addSchema(cmd, &sch, false)
	}
}

func (b *Builder) build(parent *cobra.Command, name string, cmdMap *CommandsMap) {
	if parent == nil || cmdMap == nil {
		return
	}

	if cmdMap.Command == nil {
		// and if it's a empty command, add a help
		cmdMap.Command = &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("List of '%s' commands.", name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Help()
			},
		}
	}
	// Avoid to add parent as its own child. Only should happen with root builder
	if parent != cmdMap.Command {
		if cmdMap.Command.GroupID != "" {
			if !hasGroupID(parent, cmdMap.Command.GroupID) {
				parent.AddGroup(&cobra.Group{
					Title: toTitle(cmdMap.Command.GroupID),
					ID:    cmdMap.Command.GroupID,
				})
			}
		}
		parent.AddCommand(cmdMap.Command)
	}

	if cmdMap.SubCommands == nil {
		// No more subcommands, return
		return
	}

	// Iterate over all subcommands
	for k, subCmd := range cmdMap.SubCommands {
		b.build(cmdMap.Command, k, subCmd)
	}
}

func hasGroupID(parent *cobra.Command, id string) bool {
	if parent == nil {
		return false
	}
	for _, g := range parent.Groups() {
		if g.ID == id {
			return true
		}
	}
	return false
}

func toTitle(s string) string {
	if len(s) == 0 {
		return s
	}
	firstRune := []rune(s)[0]
	upperFirstRune := unicode.ToUpper(firstRune)
	return string(upperFirstRune) + s[1:]
}

func (pc *CommandsMap) addSubCommand(path []string, newCmd *cobra.Command) error {
	if len(path) == 0 {
		pc.Command = getLatestCommand(pc.Command, newCmd)
		return nil
	}

	if pc.SubCommands == nil {
		pc.SubCommands = map[string]*CommandsMap{}
	}

	key := path[0]
	if pc.SubCommands[key] == nil {
		pc.SubCommands[key] = &CommandsMap{}
	}
	return pc.SubCommands[key].addSubCommand(path[1:], newCmd)
}

func (pc *CommandsMap) getCommand(path []string) *cobra.Command {
	// Return a copy of the command in that path if exists, any other case return nil
	if path == nil {
		return nil
	}

	if len(path) == 0 {
		return copyCommand(pc.Command)
	}

	if pc.SubCommands == nil {
		return nil
	}

	key := path[0]
	_, ok := pc.SubCommands[key]
	if !ok {
		return nil
	}

	return pc.SubCommands[key].getCommand(path[1:])
}

func copyCommand(cmd *cobra.Command) *cobra.Command {
	if cmd == nil {
		return nil
	}
	cp := *cmd

	cp.ResetFlags()
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		fcp := *f
		cp.Flags().AddFlag(&fcp)
	})

	return &cp
}

func getLatestCommand(cmd1, cmd2 *cobra.Command) *cobra.Command {
	if cmd1 == nil {
		return cmd2
	}

	if cmd2 == nil {
		return cmd1
	}

	cmd1Version, exists := cmd1.Annotations["version"]
	if !exists {
		return cmd2
	}
	cmd2Version, exists := cmd2.Annotations["version"]
	if !exists {
		return cmd1
	}

	if compareVersion(cmd1Version, cmd2Version) {
		return cmd1
	}

	return cmd2
}

func compareVersion(v1, v2 string) bool {
	// return true if v1 is newer than v2
	v1data, err := newVersionData(v1)
	if err != nil {
		return false
	}
	v2data, err := newVersionData(v2)
	if err != nil {
		return true
	}

	if v1data.IsBeta && !v2data.IsBeta {
		return false
	}

	if !v1data.IsBeta && v2data.IsBeta {
		return true
	}

	return v1data.Number >= v2data.Number
}

type versionData struct {
	IsBeta bool
	Number int
}

func newVersionData(data string) (*versionData, error) {
	isBeta := strings.Contains(data, "beta")
	ndxStart := 1
	if isBeta {
		ndxStart = 5
	}

	n, err := strconv.Atoi(data[ndxStart:])
	if err != nil {
		return nil, err
	}

	return &versionData{
		IsBeta: isBeta,
		Number: n,
	}, nil
}

// executePangeaRequest returns a function that executes a command against Pangea API, using the
// values of the flags as the body of the request.
func (b *Builder) executePangeaRequest(config *pangea.Config, svc, pathAPI string) func(cmd *cobra.Command, args []string) error {
	getClient := func(cmd *cobra.Command) (*pangea.Client, error) {
		profile, err := cmd.Flags().GetString("cli-profile")
		if err != nil {
			return nil, err
		}
		if profile != "" {
			config.Token, config.Domain, err = cli.GetProfileTokenAndDomain(profile, svc)
			if err != nil {
				return nil, err
			}
		}

		return pangea.NewClient(svc, config), nil
	}

	return func(cmd *cobra.Command, args []string) error {
		data := map[string]any{}
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if !f.Changed {
				return
			}

			// `cli-profile` flag should be ignored here
			if f.Name == "cli-profile" {
				return
			}

			sv := f.Value.String()
			// Check if it's a file
			if strings.HasPrefix(sv, "@") {
				filename := strings.ReplaceAll(sv, "@", "")
				file, err := os.Open(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", filename)
					return
				}
				defer file.Close()
				content, err := io.ReadAll(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to read file: %s\n", filename)
					return
				}
				data[f.Name] = string(content)
				return
			}

			if v, ok := f.Value.(PangeaFlag); ok {
				data[f.Name] = v.Get()
			} else {
				data[f.Name] = f.Value.String()
			}
		})

		client, err := getClient(cmd)
		if err != nil {
			return err
		}

		url, err := client.GetURL(pathAPI)
		if err != nil {
			return err
		}

		req, err := client.NewRequest("POST", url, data)
		if err != nil {
			return err
		}

		ctx := context.Background()
		var respData map[string]any
		_, err = client.Do(ctx, req, &respData, true)
		if err != nil {
			return err
		}

		cli.PrettyPrint(respData)

		b, err := json.Marshal(respData)
		if err == nil {
			if cmd.Annotations == nil {
				cmd.Annotations = map[string]string{}
			}
			cmd.Annotations["pangea_response"] = string(b)
		}

		return nil
	}
}
