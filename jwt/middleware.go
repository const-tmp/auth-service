package jwt

import (
	"auth/access"
	"auth/authz"
	"auth/logger"
	"context"
	"errors"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
)

var (
	Unauthorized = errors.New("unauthorized")
	l            = logger.New("[ jwt middleware ]\t")
)

func ValidatorFactory(a authz.Authorizer, permissions ...string) func(acc access.Access) bool {
	return func(acc access.Access) bool {
		return a.Access(permissions...).Check(acc)
	}
}

func ValidatorMiddleware(a authz.Authorizer, permissions ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			claims := ctx.Value(jwt.JWTClaimsContextKey).(*AccessClaims)
			l.Printf("claims: %+v", claims)
			if ValidatorFactory(a, permissions...)(claims.Access) {
				l.Println("authorized")
				return next(ctx, request)
			}
			l.Println("unauthorized")
			return nil, Unauthorized
		}
	}
}

func Middleware(j Service, a authz.Authorizer, permissions ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return jwt.NewParser(
			func(token *stdjwt.Token) (interface{}, error) {
				return j.PublicKey(), nil
			},
			stdjwt.SigningMethodES256,
			AccessClaimsFactory,
		)(ValidatorMiddleware(a, permissions...)(next))
	}
}
