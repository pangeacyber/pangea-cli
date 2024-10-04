// FlagMap class implements a custom Cobra map flag

package builder

import (
	"fmt"
	"strconv"
)

type FlagBool struct {
	value *bool
}

func (f *FlagBool) Type() string {
	return "boolean"
}

func (f *FlagBool) String() string {
	if f.value == nil {
		return ""
	}
	return fmt.Sprint(*f.value)
}

func (f *FlagBool) Set(in string) error {
	v, err := strconv.ParseBool(in)
	if err != nil {
		return err
	}
	f.value = &v
	return nil
}

func (f *FlagBool) Get() any {
	return f.value
}

// Ensure FlagBool implements PangeaFlag
var _ PangeaFlag = (*FlagBool)(nil)
