package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	s := New(
		logger.New("[ jwt ]\t"),
		jwt.SigningMethodES256,
		[]string{jwt.SigningMethodES256.Name},
		func() jwt.Claims { return &AccessClaims{} },
		key,
	)

	at, err := s.AccessToken(1, 1, 1, []string{"test", "test2"}, 1*time.Minute)
	require.NoError(t, err)

	t.Log(at)

	claims, err := s.VerifyAccessToken(at)
	require.NoError(t, err)

	t.Log(claims)
}

func TestJWTECDSA(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	token := jwt.New(jwt.SigningMethodES256)
	s, err := token.SignedString(key)
	require.NoError(t, err)
	t.Log(s)
	t.Log(key)

	token, err = jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("unexpected token signing method %s", token.Method.Alg())
		}
		return &key.PublicKey, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
}

func TestJWTHMAC(t *testing.T) {
	//key := make([]byte, 32)
	//_, err := rand.Reader.Read(key)
	//require.NoError(t, err)
	key := []byte("your-256-bit-secret")
	t.Log(key)
	token := jwt.New(jwt.SigningMethodHS256)
	s, err := token.SignedString(key)
	require.NoError(t, err)
	t.Log(s)
	t.Log(string(key))

	token, err = jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected token signing method %s", token.Method.Alg())
		}
		return key, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
}
