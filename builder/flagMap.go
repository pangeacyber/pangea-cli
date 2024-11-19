// FlagMap class implements a custom Cobra map flag

package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pangeacyber/pangea-cli/v2/cli"
)

type FlagMap struct {
	values map[string]any
}

func (f *FlagMap) Type() string {
	return "map"
}

func (f *FlagMap) String() string {
	b, err := json.Marshal(f.values)
	if err != nil {
		return ""
	}
	return string(b)
}

func (f *FlagMap) Set(in string) error {
	pairs, err := cli.ReadAsCSV(in)
	if err != nil {
		return err
	}
	if f.values == nil {
		f.values = make(map[string]any)
	}
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return fmt.Errorf("invalid key-value pair: %s", pair)
		}
		f.values[kv[0]] = kv[1]
	}
	return nil
}

func (f *FlagMap) Get() any {
	return f.values
}

// Append adds the specified value to the end of the flag.
func (f *FlagMap) Append(key string, value any) error {
	if f.values == nil {
		f.values = make(map[string]any)
	}
	f.values[key] = value
	return nil
}

// Replace will fully overwrite any data currently in the flag.
func (f *FlagMap) Replace(nv map[string]any) error {
	f.values = nv
	return nil
}

// GetMap returns the flag value
func (f *FlagMap) GetMap() map[string]any {
	return f.values
}

// Ensure FlagMap implements PangeaFlag
var _ PangeaFlag = (*FlagMap)(nil)
