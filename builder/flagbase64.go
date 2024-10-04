// FlagEnum class implements a custom Cobra string flag with a predefined set of values.

package builder

import (
	"encoding/base64"
)

type FlagBase64 struct {
	value string
}

func (f *FlagBase64) Type() string {
	return "string"
}

func (f *FlagBase64) String() string {
	return f.value
}

func (f *FlagBase64) Set(val string) error {
	f.value = base64.StdEncoding.EncodeToString([]byte(val))
	return nil
}

func (f *FlagBase64) Get() any {
	return f.value
}

// Ensure FlagBase64 implements PangeaFlag
var _ PangeaFlag = (*FlagBase64)(nil)
