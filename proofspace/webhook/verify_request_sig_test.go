package webhook_test

import (
	. "atala-go/proofspace/webhook"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_verifyRequestSignature(t *testing.T) {
	rsakey := newTestKey(t)
	publicPEM := encodePublicPEM(rsakey)
	tokenString := newTokenString(t, rsakey)
	headers := map[string][]string{
		AuthHeaderName: {"Bearer " + tokenString},
	}
	err := VerifyRequestSignature(headers, string(publicPEM))
	if err != nil {
		t.Errorf("error verifying request signature: %++v", err)
	}
	assert.Nil(t, err)
}

func newTokenString(t *testing.T, rsakey *rsa.PrivateKey) string {
	authToken := jwt.New(jwt.SigningMethodRS256)
	tokenString, err := authToken.SignedString(rsakey)
	if err != nil {
		t.Errorf("error signing mock jwt auth token: %++v", err)
	}
	return tokenString
}
func newTestKey(t *testing.T, bitSize_opt ...int) *rsa.PrivateKey {
	bitSize := 4096
	if len(bitSize_opt) > 0 {
		bitSize = bitSize_opt[0]
	}
	rsakey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		t.Errorf("error creating mock key pair: %++v", err)
	}
	return rsakey
}

func encodePublicPEM(rsakey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&rsakey.PublicKey),
	})
}
