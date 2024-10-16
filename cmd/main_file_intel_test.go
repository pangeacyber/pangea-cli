package main_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileIntelReputation(t *testing.T) {
	// IP Check
	hash := "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e"
	r := run("file-intel", "v2", "/reputation", "--hash_type", "sha256", "--provider", "reversinglabs", "--hashes", hash)

	data := r["data"].(map[string]any)
	assert.NotEmpty(t, data)
	hashData := data[hash].(map[string]any)
	assert.NotEmpty(t, hashData)
	category := hashData["category"].([]any)
	assert.Equal(t, len(category), 1)
	score := hashData["score"].(float64)
	assert.Equal(t, int(score), 100)
	verdict := hashData["verdict"].(string)
	assert.Equal(t, verdict, "malicious")
}

func TestFileIntelReputationBulk(t *testing.T) {
	// IP Check
	hash1 := "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e"
	hash2 := "178e2b8a4162372cd9344b81793cbf74a9513a002eda3324e6331243f3137a63"
	hashes := fmt.Sprintf("%s,%s", hash1, hash2)

	r := run("file-intel", "v2", "/reputation", "--hash_type", "sha256", "--provider", "reversinglabs", "--hashes", hashes)

	data := r["data"].(map[string]any)
	assert.NotEmpty(t, data)

	hashData := data[hash1].(map[string]any)
	assert.NotEmpty(t, hashData)
	category := hashData["category"].([]any)
	assert.Equal(t, len(category), 1)
	score := hashData["score"].(float64)
	assert.Equal(t, int(score), 100)
	verdict := hashData["verdict"].(string)
	assert.Equal(t, verdict, "malicious")
}

func TestFileIntelReputationPattern(t *testing.T) {
	r := run("file-intel", "v2", "/reputation", "file-pattern", "./testdata/*.txt,./testdata/*.json", "--hash_type", "sha256", "--provider", "reversinglabs")

	data := r["data"].(map[string]any)
	assert.NotEmpty(t, data)
}
