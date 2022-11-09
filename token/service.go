package jwt

import (
	"auth/logger"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
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
	redis                *redis.Client
	validMethods         []string
	logger               *log.Logger
}

var InvalidToken = fmt.Errorf("invalid token")

func NewService(redis *redis.Client) (*Service, error) {
	var privateKey = make([]byte, pkSize)
	_, err := rand.Read(privateKey)
	if err != nil {
		return nil, err
	}

	return &Service{
		privateKey:           privateKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		redis:                redis,
		validMethods:         []string{"HS256"},
		logger:               logger.New("jwt service\t"),
	}, nil
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID   uint64 `json:"user_id"`
	Access   uint64 `json:"access"`
	Username string `json:"username"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}

func (s *Service) NewSession(user *models.User) (string, string, error) {
	at, rt, sid, err := s.generateTokens(user)
	if err != nil {
		return "", "", err
	}
	r := s.redis.Set(context.TODO(), sid, *user.ID, s.refreshTokenDuration)
	_, err = r.Result()
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *Service) Logout(sessionID string) (int64, error) {
	return s.redis.GetDel(context.TODO(), sessionID).Int64()
}

func (s *Service) generateTokens(user *models.User) (accessToken, refreshToken, sessionID string, err error) {
	now := time.Now()
	sid, err := uuid.NewV4()
	if err != nil {
		return
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		UserID:   *user.ID,
		Access:   uint64(user.Access),
		Username: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(s.accessTokenDuration)},
			ID:        sid.String(),
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

	sessionID = sid.String()
	accessToken = accessTokenString
	refreshToken = refreshTokenString
	return
}

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
		return nil, nil, false, err
	}

	if err != nil {
		log.Println("err != nil", err)
		return nil, nil, false, err
	}

	if err := t.Claims.Valid(); err != nil {
		log.Println("validation error:", err)
		return nil, nil, false, err
	}

	if claims.UserID == 0 {
		log.Println("claims.UserID == 0", err)
		return nil, nil, false, err
	}
	userID, err := s.redis.Get(context.TODO(), claims.ID).Int64()
	if err != nil {
		log.Println("redis err != nil:", err)
		return nil, nil, false, err
	}
	if uint64(userID) != claims.UserID {
		log.Println("uint64(userID) != claims.UserID", err)
		return nil, nil, false, err
	}

	return &claims, t, true, nil
}

func (s *Service) VerifyRefreshToken(token string) (*RefreshClaims, *jwt.Token, bool, error) {
	claims := RefreshClaims{}
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
		return nil, nil, false, err
	}

	if err := t.Claims.Valid(); err != nil {
		log.Println("validation error:", err)
		return nil, nil, false, err
	}

	sid := claims.ID
	_, err = s.redis.Get(context.TODO(), sid).Result()
	if err != nil {
		return nil, nil, false, err
	}

	return &claims, t, true, nil
}
