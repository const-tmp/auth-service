//go:generate protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/service.proto
//go:generate microgen -file service.go -package github.com/nullc4t/auth-service/pkg/auth -out . -pb-go proto/service.pb.go

package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/nullc4t/auth-service/pkg/access"
	"github.com/nullc4t/auth-service/pkg/ec"
	"github.com/nullc4t/auth-service/pkg/jwt"
	"github.com/nullc4t/auth-service/pkg/mgmt"
	password2 "github.com/nullc4t/auth-service/pkg/password"
	"github.com/nullc4t/auth-service/pkg/types"

	"gorm.io/gorm"
	"log"
	"time"
)

// @microgen middleware, logging, http, grpc, recovering, error-logging
// @protobuf github.com/nullc4t/auth-service/pkg/auth/proto
type Service interface {
	Register(ctx context.Context, login, password, service string, accountId uint32) (ok bool, err error)
	Login(ctx context.Context, login, password, service string) (token *types.AccessToken, err error)
	PublicKey(ctx context.Context) (pub []byte, err error)
	GetPermissionsForService(ctx context.Context, name string) (permissions []*types.Permission, err error)
}

type service struct {
	logger *log.Logger
	mgmt   mgmt.Service
	jwt    jwt.Service
}

func (s service) GetPermissionsForService(ctx context.Context, name string) (permissions []*types.Permission, err error) {
	svc, err := s.mgmt.GetService(ctx, &types.Service{Name: name})
	if err != nil {
		return nil, fmt.Errorf("get service error: %w", err)
	}
	return s.mgmt.GetFilteredPermissions(ctx, &types.Permission{ServiceID: svc.ID})
}

func New(logger *log.Logger, mgmt mgmt.Service, jwt jwt.Service) Service {
	return &service{logger: logger, mgmt: mgmt, jwt: jwt}
}

var BadCreds = errors.New("bad credentials")

func (s service) Register(ctx context.Context, login, password, service string, accountId uint32) (bool, error) {
	u, err := s.mgmt.GetUser(ctx, &types.User{Name: login})
	if err == nil {
		s.logger.Println("user", login, "exists:", u)
		return false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	u, err = s.mgmt.CreateUserWithLoginPassword(ctx, login, password)
	if err != nil {
		s.logger.Println("create user error:", err)
		return false, err
	}

	var a *types.Account

	if accountId == 0 {
		a, err = s.mgmt.CreateAccount(ctx)
	} else {
		a, err = s.mgmt.GetAccount(ctx, &types.Account{Model: types.Model{ID: accountId}})
	}
	if err != nil {
		s.logger.Printf("get or create account %d error: %s", accountId, err)
		return false, err
	}

	ok, err := s.mgmt.AttachUserToAccount(ctx, u.ID, a.ID)
	if err != nil {
		s.logger.Printf("set account %d error: %s", accountId, err)
		return ok, err
	}

	if service != "" {
		var svc *types.Service
		svc, err = s.mgmt.GetService(ctx, &types.Service{Name: service})
		if err != nil {
			s.logger.Println("get auth error:", err)
			return false, err
		}

		ok, err := s.mgmt.AttachAccountToService(ctx, svc.ID, a.ID)
		if err != nil {
			s.logger.Println("attach account to auth error:", err)
			return false, err
		}
		if !ok {
			s.logger.Println("attach account to auth not ok")
			return false, nil
		}
	}

	return true, nil
}

func (s service) Login(ctx context.Context, login, password, service string) (*types.AccessToken, error) {
	u, err := s.mgmt.GetUser(ctx, &types.User{Name: login})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Println("get user error", err)
			return nil, err
		}
		s.logger.Println("user", login, "not found")
		return nil, BadCreds
	}
	if u.ID == 0 {
		return nil, BadCreds
	}
	if !password2.CheckHash(password, u.Password) {
		return nil, BadCreds
	}

	svc, err := s.mgmt.GetService(ctx, &types.Service{Name: service})
	if err != nil {
		s.logger.Println("get auth error:", err)
		return nil, err
	}

	p, err := s.mgmt.GetUserPermissions(ctx, u.ID)
	if err != nil {
		s.logger.Println("get user permission error:", err)
		return nil, err
	}
	ac := access.Access(0)
	for _, permission := range p {
		ac |= permission.Access
	}
	s.logger.Println("user", u.ID, service, "permissions:", ac)

	at, err := s.jwt.AccessToken(u.ID, u.AccountID, ac, []string{svc.Name}, time.Minute*5)

	if err != nil {
		s.logger.Println("jwt error:", err)
		return nil, err
	}

	return &types.AccessToken{AccessToken: at}, nil
}

func (s service) PublicKey(ctx context.Context) (pub []byte, err error) {
	return ec.PublicKey2PEM(s.jwt.PublicKey())
}
