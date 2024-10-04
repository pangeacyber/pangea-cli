// FlagMap class implements a custom Cobra map flag

package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pangeacyber/pangea-cli-internal/cli"
)

type FlagAny struct {
	values any
}

func (f *FlagAny) Type() string {
	return "any"
}

func (f *FlagAny) String() string {
	b, err := json.Marshal(f.values)
	if err != nil {
		return ""
	}
	return string(b)
}

func (f *FlagAny) Set(in string) error {
	pairs, err := cli.ReadAsCSV(in)
	if err != nil {
		f.values = in
		return nil
	}

	m := make(map[string]any)
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return fmt.Errorf("invalid key-value pair: %s", pair)
		}
		m[kv[0]] = kv[1]
	}
	f.values = m
	return nil
}

func (f *FlagAny) Get() any {
	return f.values
}

// Replace will fully overwrite any data currently in the flag.
func (f *FlagAny) Replace(nv map[string]any) error {
	f.values = nv
	return nil
}

// Ensure FlagAny implements PangeaFlag
var _ PangeaFlag = (*FlagAny)(nil)
