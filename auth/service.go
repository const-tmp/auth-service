package auth

import (
	"auth/account"
	"auth/auth/authz"
	"auth/jwt"
	"auth/mgmt"
	"auth/pkg/password"
	"auth/pkg/types"
	svcsrv "auth/service"
	"auth/user"
	"context"
	"errors"
	"github.com/nullc4ts/bitmask_authz/access"
	"gorm.io/gorm"
	"log"
	"time"
)

// @microgen middleware, logging, http, recovering, error-logging
type Service interface {
	Register(ctx context.Context, login, password, service string, accountID uint) (ok bool, err error)
	Login(ctx context.Context, login, password, service string) (at types.AccessToken, err error)
}

type service struct {
	logger  *log.Logger
	user    user.Service
	authz   authz.Service
	mgmt    mgmt.Service
	svc     svcsrv.Service
	account account.Repo
	jwt     jwt.Service
}

func New(logger *log.Logger, user user.Service, authz authz.Service, mgmt mgmt.Service, svc svcsrv.Service, account account.Repo, jwt jwt.Service) Service {
	return &service{logger: logger, user: user, authz: authz, mgmt: mgmt, svc: svc, account: account, jwt: jwt}
}

func (s service) Register(ctx context.Context, login, password, service string, accountID uint) (bool, error) {
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

	var a types.Account

	if accountID == 0 {
		a, err = s.account.Create(ctx)
	} else {
		a, err = s.account.Get(ctx, types.Account{Model: gorm.Model{ID: accountID}})
	}
	if err != nil {
		s.logger.Printf("get or create account %d error: %s", accountID, err)
		return false, err
	}

	ok, err := s.user.SetAccount(ctx, u.ID, a.ID)
	if err != nil {
		s.logger.Printf("set account %d error: %s", accountID, err)
		return ok, err
	}

	if service != "" {
		var svc types.Service
		svc, err = s.svc.Get(ctx, types.Service{Name: service})
		if err != nil {
			s.logger.Println("get service error:", err)
			return false, err
		}

		ok, err := s.mgmt.AttachAccountToService(ctx, svc.ID, a.ID)
		if err != nil {
			s.logger.Println("attach account to service error:", err)
			return false, err
		}
		if !ok {
			s.logger.Println("attach account to service not ok")
			return false, nil
		}
	}

	return true, nil
}

var BadCreds = errors.New("bad credentials")

func (s service) Login(ctx context.Context, login, pass, service string) (types.AccessToken, error) {
	u, err := s.user.Get(ctx, types.User{Name: login})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Println("get user error", err)
			return types.AccessToken{}, err
		}
		s.logger.Println("user", login, "not found")
		return types.AccessToken{}, BadCreds
	}
	if u.ID == 0 {
		return types.AccessToken{}, BadCreds
	}
	if !password.CheckHash(pass, u.Password) {
		return types.AccessToken{}, BadCreds
	}

	svc, err := s.svc.Get(ctx, types.Service{Name: service})
	if err != nil {
		s.logger.Println("get service error:", err)
		return types.AccessToken{}, err
	}

	p, err := s.authz.GetUserPermissions(ctx, types.Permission{ServiceID: svc.ID}, u.ID)
	if err != nil {
		s.logger.Println("get user permission error:", err)
		return types.AccessToken{}, err
	}
	ac := access.Access(0)
	for _, permission := range p {
		ac |= permission.Access
	}
	s.logger.Println("user", u.ID, service, "permissions:", ac)

	at, err := s.jwt.AccessToken(u.ID, u.AccountID, ac, []string{svc.Name}, time.Minute*5)

	if err != nil {
		s.logger.Println("jwt error:", err)
		return types.AccessToken{}, err
	}

	return types.AccessToken{AccessToken: at}, nil
}
