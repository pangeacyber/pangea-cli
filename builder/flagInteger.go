// FlagMap class implements a custom Cobra map flag

package builder

import (
	"fmt"
	"strconv"
)

type FlagInteger struct {
	value *int
}

func (f *FlagInteger) Type() string {
	return "integer"
}

func (f *FlagInteger) String() string {
	if f.value == nil {
		return ""
	}
	return fmt.Sprint(*f.value)
}

func (f *FlagInteger) Set(in string) error {
	v, err := strconv.Atoi(in)
	if err != nil {
		return err
	}
	f.value = &v
	return nil
}

func (f *FlagInteger) Get() any {
	return f.value
}

// Ensure FlagInteger implements PangeaFlag
var _ PangeaFlag = (*FlagInteger)(nil)
