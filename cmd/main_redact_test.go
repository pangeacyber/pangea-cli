package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactText(t *testing.T) {
	redacted := "My Phone number is <PHONE_NUMBER>"
	input := "My Phone number is 415-867-5309"

	r := run("redact", "v1", "/redact", "--text", input)
	assert.Equal(t, int(r["count"].(float64)), 1)
	assert.Equal(t, r["redacted_text"].(string), redacted)
}

func TestRedactStructured(t *testing.T) {
	input := "phone=415-867-5309"
	redacted := map[string]any{"phone": "<PHONE_NUMBER>"}

	r := run("redact", "v1", "/redact_structured", "--data", input)
	assert.Equal(t, int(r["count"].(float64)), 1)
	assert.Equal(t, r["redacted_data"].(map[string]any), redacted)
}
