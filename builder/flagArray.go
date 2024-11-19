// FlagArray class implements a custom Cobra string array flag.

package builder

import (
	"encoding/json"

	"github.com/pangeacyber/pangea-cli/v2/cli"
)

type FlagArray struct {
	values []string
}

func (f *FlagArray) Type() string {
	return "array"
}

func (f *FlagArray) String() string {
	b, err := json.Marshal(f.values)
	if err != nil {
		return ""
	}
	return string(b)
}

func (f *FlagArray) Set(in string) error {
	v, err := cli.ReadAsCSV(in)
	if err != nil {
		return err
	}
	if f.values == nil {
		f.values = v
	} else {
		f.values = append(f.values, v...)
	}
	return nil
}

func (f *FlagArray) Get() any {
	return f.values
}

// Append adds the specified value to the end of the flag value list.
func (f *FlagArray) Append(nv string) error {
	f.values = append(f.values, nv)
	return nil
}

// Replace will fully overwrite any data currently in the flag value list.
func (f *FlagArray) Replace(nv []string) error {
	f.values = nv
	return nil
}

// GetSlice returns the flag value list as an array of strings.
func (f *FlagArray) GetSlice() []string {
	return f.values
}

// Ensure FlagArray implements PangeaFlag
var _ PangeaFlag = (*FlagArray)(nil)
