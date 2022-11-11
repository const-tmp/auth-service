package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nullc4t/auth-service/pkg/access"
	"time"
)

type ClaimsFactory func() jwt.Claims

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID    uint32        `json:"user_id"`
	AccountID uint32        `json:"account_id"`
	Access    access.Access `json:"access"`
}

func (c AccessClaims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return fmt.Errorf("RegisteredClaims validation error: %w", err)
	}
	if c.UserID == 0 {
		return fmt.Errorf("AccessClaims validation error: UserID = %d", c.UserID)
	}
	return nil
}

func AccessClaimsFactory() jwt.Claims { return &AccessClaims{} }

func NewAccessClaims(userID, accID uint32, acc access.Access, aud []string, duration time.Duration) AccessClaims {
	now := time.Now()
	//sid, _ := uuid.NewV4()
	return AccessClaims{
		UserID:    userID,
		AccountID: accID,
		Access:    acc,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(duration)},
			Audience:  aud,
			//ID:        sid.String(),
		},
	}
}
