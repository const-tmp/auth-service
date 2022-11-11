package jwt

import (
	"context"
	"errors"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
	"github.com/nullc4t/auth-service/pkg/access"
	"github.com/nullc4t/auth-service/pkg/logger"
)

var (
	Unauthorized = errors.New("unauthorized")
	l            = logger.New("[ jwt middleware ]\t")
)

func ValidatorFactory(a access.Helper, permissions ...string) func(acc access.Access) bool {
	return func(acc access.Access) bool {
		return a.Access(permissions...).Check(acc)
	}
}

func AccessCheckMiddleware(a access.Helper, permissions ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			claims := ctx.Value(kitjwt.JWTClaimsContextKey).(*AccessClaims)
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

func Middleware(j Service, a access.Helper, permissions ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return kitjwt.NewParser(
			func(token *stdjwt.Token) (interface{}, error) {
				return j.PublicKey(), nil
			},
			j.SigningMethod(),
			kitjwt.ClaimsFactory(j.ClaimsFactory()),
		)(AccessCheckMiddleware(a, permissions...)(next))
	}
}

func KitAdapterMiddleware(j Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return kitjwt.NewParser(
			func(token *stdjwt.Token) (interface{}, error) {
				return j.PublicKey(), nil
			},
			j.SigningMethod(),
			kitjwt.ClaimsFactory(j.ClaimsFactory()),
		)(next)
	}
}
