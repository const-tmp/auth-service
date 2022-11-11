package ec

import (
	"auth/pkg/logger"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var l = logger.New("[ EC ]\t")

func PrivateKey2PEM(k *ecdsa.PrivateKey) ([]byte, error) {
	derBytes, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return nil, err
	}
	p := pem.Block{
		Type:    "EC PRIVATE KEY",
		Headers: nil,
		Bytes:   derBytes,
	}
	return pem.EncodeToMemory(&p), nil
}

func PEM2PrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	pemBlock, rest := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("PEM block is nil")
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("rest is not empty: %s", string(rest))
	}
	return x509.ParseECPrivateKey(pemBlock.Bytes)
}

func PublicKey2PEM(k *ecdsa.PublicKey) ([]byte, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(k)
	if err != nil {
		return nil, err
	}
	p := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   derBytes,
	}
	return pem.EncodeToMemory(&p), nil
}

func PEM2PublicKey(pemBytes []byte) (*ecdsa.PublicKey, error) {
	pemBlock, rest := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("PEM block is nil")
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("rest is not empty: %s", string(rest))
	}
	pkix, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pkix.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not OK")
	}
	return pub, nil
}
