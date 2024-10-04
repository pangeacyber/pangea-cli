package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbargoIPcheck(t *testing.T) {
	// IP Check
	r := run("embargo", "v1", "/ip/check", "--ip", "213.24.238.26")
	assert.Equal(t, int(r["count"].(float64)), 1)
	sanctions := r["sanctions"].([]any)
	assert.Equal(t, len(sanctions), 1)
	for _, item := range sanctions {
		sanction := item.(map[string]any)
		assert.NotEmpty(t, sanction["list_name"].(string))
		assert.NotEmpty(t, sanction["embargoed_country_name"].(string))
		assert.NotEmpty(t, sanction["embargoed_country_iso_code"].(string))
		assert.NotEmpty(t, sanction["issuing_country"].(string))
	}
}

func TestEmbargoISOcheck(t *testing.T) {
	// ISO Check
	r := run("embargo", "v1", "/iso/check", "--iso_code", "CU")
	assert.Equal(t, int(r["count"].(float64)), 1)
	sanctions := r["sanctions"].([]any)
	assert.Equal(t, len(sanctions), 1)
	for _, item := range sanctions {
		sanction := item.(map[string]any)
		assert.NotEmpty(t, sanction["list_name"].(string))
		assert.NotEmpty(t, sanction["embargoed_country_name"].(string))
		assert.NotEmpty(t, sanction["embargoed_country_iso_code"].(string))
		assert.NotEmpty(t, sanction["issuing_country"].(string))
	}
}
