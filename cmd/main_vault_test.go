package main_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsymmetricKeyCycle(t *testing.T) {
	// Create key
	r := run("vault", "v1", "/key/generate", "--type", "asymmetric_key", "--purpose", "signing", "--algorithm", "ED25519")
	assert.NotEmpty(t, r["id"])
	id := r["id"].(string)
	assert.NotEmpty(t, id)
	assert.Equal(t, int(r["version"].(float64)), 1)

	// Rotate key
	r = run("vault", "v1", "/key/rotate", "--id", id, "--rotation_state", "deactivated")
	assert.Equal(t, id, r["id"].(string))
	assert.Equal(t, int(r["version"].(float64)), 2)

	// Sign
	message := "aaaabbbbcccc"
	r = run("vault", "v1", "/key/sign", "--id", id, "--message", message)
	signature := r["signature"].(string)

	// Verify
	r = run("vault", "v1", "/key/verify", "--id", id, "--message", message, "--signature", signature)
	assert.True(t, r["valid_signature"].(bool))

	r = run("vault", "v1", "/key/verify", "--id", id, "--message", "abcd", "--signature", signature)
	assert.False(t, r["valid_signature"].(bool))
}

func TestSymmetricKeyCycle(t *testing.T) {
	// Create key
	r := run("vault", "v1", "/key/generate", "--type", "symmetric_key", "--purpose", "encryption", "--algorithm", "AES-CFB-128")
	assert.NotEmpty(t, r["id"])
	id := r["id"].(string)
	assert.NotEmpty(t, id)
	assert.Equal(t, int(r["version"].(float64)), 1)

	// Rotate key
	r = run("vault", "v1", "/key/rotate", "--id", id, "--rotation_state", "deactivated")
	assert.Equal(t, id, r["id"].(string))
	assert.Equal(t, int(r["version"].(float64)), 2)

	// Encrypt
	plainText := "aaaabbbbcccc"
	r = run("vault", "v1", "/key/encrypt", "--id", id, "--plain_text", plainText)
	cipherText := r["cipher_text"].(string)

	// Decrypt
	r = run("vault", "v1", "/key/decrypt", "--id", id, "--cipher_text", cipherText)
	decPlainText := r["plain_text"].(string)
	assert.Equal(t, plainText, decPlainText)
}

func TestFPEKeyCycle(t *testing.T) {
	// Create key
	r := run("vault", "v1", "/key/generate", "--type", "symmetric_key", "--purpose", "fpe", "--algorithm", "AES-FF3-1-128-BETA")
	assert.NotEmpty(t, r["id"])
	id := r["id"].(string)
	assert.NotEmpty(t, id)
	assert.Equal(t, int(r["version"].(float64)), 1)

	// Rotate key
	r = run("vault", "v1", "/key/rotate", "--id", id, "--rotation_state", "deactivated")
	assert.Equal(t, id, r["id"].(string))
	assert.Equal(t, int(r["version"].(float64)), 2)

	// Encrypt
	plainText := "123-4567-8901"
	r = run("vault", "v1", "/key/encrypt/transform", "--id", id, "--plain_text", plainText, "--alphabet", "numeric")
	cipherText := r["cipher_text"].(string)
	tweak := r["tweak"].(string)

	// Decrypt
	r = run("vault", "v1", "/key/decrypt/transform", "--id", id, "--cipher_text", cipherText, "--tweak", tweak, "--alphabet", "numeric")
	decPlainText := r["plain_text"].(string)
	assert.Equal(t, plainText, decPlainText)
}

func TestSecretCycle(t *testing.T) {
	secret1 := "secret1"
	secret2 := "secret2"
	name := fmt.Sprintf("secret_%d", time.Now().UnixMilli())
	metadata := "key:value1,key2:2"
	tags := "tag1,tag2"

	r := run("vault", "v1", "/secret/store", "--type", "secret", "--secret", secret1, "--name", name, "--folder", "/cli/secrets", "--tags", tags, "--metadata", metadata)

	id := r["id"].(string)
	assert.NotEmpty(t, id)
	assert.Equal(t, int(r["version"].(float64)), 1)
	assert.Equal(t, r["secret"].(string), secret1)
	assert.Equal(t, r["type"].(string), "secret")

	r = run("vault", "v1", "/secret/rotate", "--id", id, "--secret", secret2, "--rotation_state", "deactivated")
	id = r["id"].(string)
	assert.NotEmpty(t, id)
	assert.Equal(t, int(r["version"].(float64)), 2)
	assert.Equal(t, r["secret"].(string), secret2)
	assert.Equal(t, r["type"].(string), "secret")

	r = run("vault", "v1", "/get", "--id", id)
	assert.Equal(t, r["id"].(string), id)
	assert.Equal(t, r["type"].(string), "secret")
}

func TestFolder(t *testing.T) {
	name := fmt.Sprintf("folder_%d", time.Now().UnixMilli())
	r := run("vault", "v1", "/folder/create", "--name", name, "--folder", "/")
	id := r["id"].(string)
	assert.NotEmpty(t, id)
}

func TestList(t *testing.T) {
	r := run("vault", "v1", "/list", "--order", "asc", "--order_by", "id")
	items := r["items"].([]any)
	assert.True(t, len(items) > 0)
	for _, item := range items {
		o := item.(map[string]any)
		assert.NotEmpty(t, o["id"].(string))
	}
}

func TestStoreUsingFile(t *testing.T) {
	// Create key
	r := run("vault", "v1", "/key/store", "--type", "asymmetric_key", "--purpose", "signing", "--algorithm", "ED25519", "--public_key", "@testdata/key.pub", "--private_key", "@testdata/key.pem")
	assert.NotEmpty(t, r["id"])
	id := r["id"].(string)
	assert.NotEmpty(t, id)
	assert.Equal(t, int(r["version"].(float64)), 1)
}
