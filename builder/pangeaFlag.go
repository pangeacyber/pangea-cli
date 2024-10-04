package builder

import "github.com/spf13/pflag"

type PangeaFlag interface {
	pflag.Value

	Get() any
}
