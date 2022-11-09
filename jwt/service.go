package jwt

import (
	"auth/logger"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nullc4ts/bitmask_authz/access"
	"log"
	"time"
)

const (
	pkSize               = 32
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type Service struct {
	privateKey           []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	//redis                *redis.Client
	validMethods []string
	logger       *log.Logger
	//ecdsa        *ecdsa.PrivateKey
}

func NewService(atd, rtd time.Duration, pk []byte) Service {
	return Service{
		privateKey: pk,
		//ecdsa:                key,
		accessTokenDuration:  atd,
		refreshTokenDuration: rtd,
		validMethods:         []string{"HS256"},
		logger:               logger.New("[ jwt service ]\t"),
	}
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID    uint          `json:"user_id"`
	AccountID uint          `json:"account_id"`
	Access    access.Access `json:"access"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}

//func (s *Service) Logout(sessionID string) (int64, error) {
//	panic("unimplemented")
//	//return 0, nil
//	//return s.redis.GetDel(context.TODO(), sessionID).Int64()
//}

func (s *Service) GenerateTokens(userID, accID uint, acc access.Access, aud string) (accessToken, refreshToken string, err error) {
	now := time.Now()
	sid, err := uuid.NewV4()
	if err != nil {
		return
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		UserID:    userID,
		AccountID: accID,
		Access:    acc,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(s.accessTokenDuration)},
			ID:        sid.String(),
			Audience:  []string{aud},
		},
	})
	accessTokenString, err := at.SignedString(s.privateKey)
	if err != nil {
		return
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(s.refreshTokenDuration)},
		ID:        sid.String(),
	})
	refreshTokenString, err := rt.SignedString(s.privateKey)
	if err != nil {
		return
	}

	accessToken = accessTokenString
	refreshToken = refreshTokenString
	return
}

var (
	// ErrTokenContextMissing denotes a token was not passed into the parsing
	// middleware's context.
	ErrTokenContextMissing = errors.New("token up for parsing was not passed through the context")

	// ErrTokenInvalid denotes a token was not able to be validated.
	ErrTokenInvalid = errors.New("JWT was invalid")

	// ErrTokenExpired denotes a token's expire header (exp) has since passed.
	ErrTokenExpired = errors.New("JWT is expired")

	// ErrTokenMalformed denotes a token was not formatted as a JWT.
	ErrTokenMalformed = errors.New("JWT is malformed")

	// ErrTokenNotActive denotes a token's not before header (nbf) is in the
	// future.
	ErrTokenNotActive = errors.New("token is not valid yet")

	// ErrUnexpectedSigningMethod denotes a token was signed with an unexpected
	// signing method.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

func (s *Service) VerifyAccessToken(token string) (*AccessClaims, *jwt.Token, bool, error) {
	claims := AccessClaims{}
	t, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return s.privateKey, nil
		},
		jwt.WithValidMethods(s.validMethods),
	)
	if err != nil {
		s.logger.Println("verify access token error:", err)
		return nil, nil, false, err
	}

	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			switch {
			case e.Errors&jwt.ValidationErrorMalformed != 0:
				// Token is malformed
				return nil, nil, false, ErrTokenMalformed
			case e.Errors&jwt.ValidationErrorExpired != 0:
				// Token is expired
				return nil, nil, false, ErrTokenExpired
			case e.Errors&jwt.ValidationErrorNotValidYet != 0:
				// Token is not active yet
				return nil, nil, false, ErrTokenNotActive
			case e.Inner != nil:
				// report e.Inner
				return nil, nil, false, e.Inner
			}
			// We have a ValidationError but have no specific Go kit error for it.
			// Fall through to return original error.
		}
		return nil, nil, false, err
	}

	if !t.Valid {
		return nil, nil, false, ErrTokenInvalid
	}

	if err := t.Claims.Valid(); err != nil {
		s.logger.Println("validation error:", err)
		return nil, nil, false, err
	}

	if claims.UserID == 0 {
		s.logger.Println("claims.UserID == 0", err)
		return nil, nil, false, err
	}
	//userID, err := s.redis.Get(context.TODO(), claims.ID).Int64()
	//if err != nil {
	//	log.Println("redis err != nil:", err)
	//	return nil, nil, false, err
	//}
	//if uint64(userID) != claims.UserID {
	//	log.Println("uint64(userID) != claims.UserID", err)
	//	return nil, nil, false, err
	//}

	return &claims, t, true, nil
}

//func (s *Service) VerifyRefreshToken(token string) (*RefreshClaims, *jwt.Token, bool, error) {
//	claims := RefreshClaims{}
//	t, err := jwt.ParseWithClaims(
//		token,
//		&claims,
//		func(token *jwt.Token) (interface{}, error) {
//			_, ok := token.Method.(*jwt.SigningMethodHMAC)
//			if !ok {
//				return nil, fmt.Errorf("unexpected token signing method")
//			}
//			return s.privateKey, nil
//		},
//		jwt.WithValidMethods(s.validMethods),
//	)
//	if err != nil {
//		return nil, nil, false, err
//	}
//
//	if err := t.Claims.Valid(); err != nil {
//		log.Println("validation error:", err)
//		return nil, nil, false, err
//	}
//
//	//sid := claims.ID
//	//_, err = s.redis.Get(context.TODO(), sid).Result()
//	//if err != nil {
//	//	return nil, nil, false, err
//	//}
//
//	return &claims, t, true, nil
//}
