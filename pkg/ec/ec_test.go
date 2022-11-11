package ec

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrivateKey(t *testing.T) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	t.Log(pk)
	require.NoError(t, err)

	pemBytes, err := PrivateKey2PEM(pk)
	require.NoError(t, err)
	t.Log(pemBytes)
	t.Log(string(pemBytes))

	pk, err = PEM2PrivateKey(pemBytes)
	require.NoError(t, err)
	t.Log(pk)
}

func TestPublicKey(t *testing.T) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	t.Log(pk)
	require.NoError(t, err)

	pemBytes, err := PublicKey2PEM(&pk.PublicKey)
	require.NoError(t, err)
	t.Log(pemBytes)
	t.Log(string(pemBytes))

	pub, err := PEM2PublicKey(pemBytes)
	require.NoError(t, err)
	t.Log(pub)
}

func TestPrivateKeyJWT(t *testing.T) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	t.Log(pk)
	require.NoError(t, err)

	pemBytes, err := PrivateKey2PEM(pk)
	require.NoError(t, err)
	t.Log(pemBytes)
	t.Log(string(pemBytes))

	pk, err = jwt.ParseECPrivateKeyFromPEM(pemBytes)
	require.NoError(t, err)
	t.Log(pk)
}

func TestPublicKeyJWT(t *testing.T) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	t.Log(pk)
	require.NoError(t, err)

	pemBytes, err := PublicKey2PEM(&pk.PublicKey)
	require.NoError(t, err)
	t.Log(pemBytes)
	t.Log(string(pemBytes))

	pub, err := jwt.ParseECPublicKeyFromPEM(pemBytes)
	require.NoError(t, err)
	t.Log(pub)
}
