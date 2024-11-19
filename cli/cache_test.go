package cli_test

import (
	"fmt"
	"testing"

	"github.com/pangeacyber/pangea-cli/v2/cli"
	"github.com/stretchr/testify/assert"
)

func TestLoadCacheData(t *testing.T) {
	cache, err := cli.LoadCacheData()
	assert.NoError(t, err)

	s, err := cli.IndentedString(cache)
	assert.NoError(t, err)

	fmt.Println(s)
}
