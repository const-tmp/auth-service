//go:generate protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/service.proto
//go:generate microgen -file service.go -package auth/auth -out . -pb-go proto/service.pb.go -main

package auth

import (
	"auth/access"
	"auth/account"
	"auth/jwt"
	"auth/mgmt"
	"auth/pkg/ec"
	password2 "auth/pkg/password"
	"auth/pkg/types"
	svcsrv "auth/service"
	"auth/user"
	"context"
	"errors"

	"gorm.io/gorm"
	"log"
	"time"
)

// @microgen middleware, logging, http, grpc, recovering, error-logging
// @protobuf auth/auth/proto
type Service interface {
	Register(ctx context.Context, login, password, service string, accountId uint32) (ok bool, err error)
	Login(ctx context.Context, login, password, service string) (token *types.AccessToken, err error)
	PublicKey(ctx context.Context) (pub []byte, err error)
}

type service struct {
	logger *log.Logger
	user   user.Service
	//authz   authz.Service
	mgmt    mgmt.Service
	svc     svcsrv.Service
	account account.Repo
	jwt     jwt.Service
}

func New(logger *log.Logger, user user.Service, mgmt mgmt.Service, svc svcsrv.Service, account account.Repo, jwt jwt.Service) Service {
	return &service{logger: logger, user: user, mgmt: mgmt, svc: svc, account: account, jwt: jwt}
}

var BadCreds = errors.New("bad credentials")

func (s service) Register(ctx context.Context, login, password, service string, accountId uint32) (bool, error) {
	u, err := s.user.Get(ctx, types.User{Name: login})
	if err == nil {
		s.logger.Println("user", login, "exists:", u)
		return false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	u, err = s.user.CreateWithLoginPassword(ctx, login, password)
	if err != nil {
		s.logger.Println("create user error:", err)
		return false, err
	}

	var a *types.Account

	if accountId == 0 {
		a, err = s.account.Create(ctx)
	} else {
		a, err = s.account.Get(ctx, &types.Account{Model: types.Model{ID: accountId}})
	}
	if err != nil {
		s.logger.Printf("get or create account %d error: %s", accountId, err)
		return false, err
	}

	ok, err := s.user.SetAccount(ctx, u.ID, a.ID)
	if err != nil {
		s.logger.Printf("set account %d error: %s", accountId, err)
		return ok, err
	}

	if service != "" {
		var svc *types.Service
		svc, err = s.svc.Get(ctx, &types.Service{Name: service})
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
	u, err := s.user.Get(ctx, types.User{Name: login})
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

	svc, err := s.svc.Get(ctx, &types.Service{Name: service})
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
