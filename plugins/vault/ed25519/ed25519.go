package ed25519

import (
	"crypto"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/pangeacyber/pangea-cli/v2/plugins/vault/common"
)

// Generates Ed25519 key pairs
func GenerateKeyPair() (pubKey ed25519.PublicKey, privKey ed25519.PrivateKey, err error) {
	rand := cryptorand.Reader

	seed := make([]byte, ed25519.SeedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, nil, fmt.Errorf("generate asymmetric key failed: %w", err)
	}

	pubKey, privKey, err = ed25519.GenerateKey(rand)
	if err != nil {
		return nil, nil, fmt.Errorf("generate asymmetric key failed: %w", err)
	}

	return pubKey, privKey, nil
}

// Encode Private Key to PKCS #8, ASN.1 DER format embedded in a PEM Block
func EncodePEMPrivateKey(privKey crypto.PrivateKey) ([]byte, error) {
	pkcs, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, common.ErrInvalidPrivateKey
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs,
	}
	return pem.EncodeToMemory(block), nil
}

// Encode Public Key to PKIX, ASN.1 DER format embedded in a PEM Block
func EncodePEMPublicKey(pubKey crypto.PublicKey) ([]byte, error) {
	pkix, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, common.ErrInvalidPublicKey
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkix,
	}
	return pem.EncodeToMemory(block), nil
}
