package utils

import (
	"fmt"
	"testing"
)

func TestReplaceHomeFolder(t *testing.T) {
	patters := []string{"~/pangea/*.json", "/tmp/*.txt", "pangea~path/*.txt", "~pangea/*.txt"}
	newPatters := replaceUserFolder(patters)

	for _, p := range newPatters {
		fmt.Println(p)
	}
}
