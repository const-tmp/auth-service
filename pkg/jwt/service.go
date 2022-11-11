package jwt

import (
	"auth/pkg/access"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type Service interface {
	AccessToken(userID, accID uint32, acc access.Access, aud []string, duration time.Duration) (string, error)
	VerifyAccessToken(token string) (jwt.Claims, error)
	Key() *ecdsa.PrivateKey
	PublicKey() *ecdsa.PublicKey
}

type service struct {
	logger              *log.Logger
	signingMethod       jwt.SigningMethod
	validSigningMethods []string
	claimsFactory       ClaimsFactory
	key                 *ecdsa.PrivateKey
}

func New(
	logger *log.Logger,
	signingMethod jwt.SigningMethod,
	validSigningMethods []string,
	claimsFactory ClaimsFactory,
	key *ecdsa.PrivateKey,
) Service {
	return &service{
		logger:              logger,
		signingMethod:       signingMethod,
		validSigningMethods: validSigningMethods,
		claimsFactory:       claimsFactory,
		key:                 key,
	}
}

func (s service) AccessToken(userID, accID uint32, acc access.Access, aud []string, duration time.Duration) (string, error) {
	t := jwt.NewWithClaims(s.signingMethod, NewAccessClaims(userID, accID, acc, aud, duration))

	pk, err := s.PrivateKeyFunc(t)
	if err != nil {
		return "", err
	}

	accessToken, err := t.SignedString(pk)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s service) VerifyAccessToken(tokenString string) (jwt.Claims, error) {
	claims := s.claimsFactory()
	token, err := jwt.ParseWithClaims(tokenString, claims, s.PublicKeyFunc, jwt.WithValidMethods(s.validSigningMethods))
	if err != nil {
		s.logger.Println("parse with claims error:", err)
		if e, ok := err.(*jwt.ValidationError); ok {
			switch {
			case e.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, errors.New("JWT is malformed")
			case e.Errors&jwt.ValidationErrorExpired != 0:
				return nil, errors.New("JWT is expired")
			case e.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, errors.New("token is not valid yet")
			case e.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, errors.New("JWT signature is invalid")
			case e.Inner != nil:
				return nil, e.Inner
			default:
				return nil, e
			}
		}
		return nil, fmt.Errorf("parse with claims error: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("JWT token is invalid")
	}

	if err := token.Claims.Valid(); err != nil {
		s.logger.Println("claims validation error:", err)
		return nil, fmt.Errorf("claims validation error: %w", err)
	}

	return claims, nil
}

func (s service) Key() *ecdsa.PrivateKey {
	return s.key
}

func (s service) PublicKey() *ecdsa.PublicKey {
	return &s.Key().PublicKey
}

func (s service) PrivateKeyFunc(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodECDSA)
	if !ok {
		return nil, fmt.Errorf("unexpected token signing method")
	}
	return s.key, nil
}

func (s service) PublicKeyFunc(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodECDSA)
	if !ok {
		return nil, fmt.Errorf("unexpected token signing method")
	}
	return &s.key.PublicKey, nil
}
