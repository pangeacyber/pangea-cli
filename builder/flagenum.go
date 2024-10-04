// FlagEnum class implements a custom Cobra string flag with a predefined set of values.

package builder

import (
	"fmt"
	"strings"
)

type FlagEnum struct {
	enum    []string
	enumMap map[string]bool
	Name    string
	value   string
}

func NewFlagEnum(name string, values []string) *FlagEnum {
	f := FlagEnum{
		Name: name,
	}
	f.AddValues(values)
	return &f
}

func (f *FlagEnum) Type() string {
	return "string (enum)"
}

func (f *FlagEnum) String() string {
	return f.value
}

func (f *FlagEnum) Set(val string) error {
	for _, v := range f.enum {
		if val == v {
			f.value = val
			return nil
		}
	}
	return fmt.Errorf("invalid value %s. Possible values: [%s]", val, strings.Join(f.enum, " "))
}

func (f *FlagEnum) Get() any {
	return f.value
}

func (f *FlagEnum) Description() string {
	return fmt.Sprintf("Possible values: [%s]", strings.Join(f.enum, " "))
}

func (f *FlagEnum) AddValues(values []string) {
	if f.enumMap == nil {
		f.enumMap = make(map[string]bool, 0)
	}
	for _, v := range values {
		v = strings.TrimSpace(v)
		if flag := f.enumMap[v]; !flag {
			f.enumMap[v] = true
			f.enum = append(f.enum, v)
		}
	}
}

func (f *FlagEnum) GetValues() []string {
	return f.enum
}

// Ensure FlagEnum implements PangeaFlag
var _ PangeaFlag = (*FlagEnum)(nil)
