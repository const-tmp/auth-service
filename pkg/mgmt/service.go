//go:generate protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/service.proto
//go:generate microgen -file service.go -package auth/mgmt -out . -pb-go proto/service.pb.go

package mgmt

import (
	"auth/pkg/access"
	"auth/pkg/account"
	"auth/pkg/permission"
	svcsrv "auth/pkg/service"
	"auth/pkg/types"
	"auth/pkg/user"
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
)

// @microgen middleware, logging, http, grpc, recovering, error-logging
// @protobuf auth/mgmt/proto
type Service interface {
	CreateUserWithLoginPassword(ctx context.Context, login, pass string) (user *types.User, err error)
	CreateUserWithTelegram(ctx context.Context, id uint64, name, userN string) (user *types.User, err error)
	GetAllUsers(ctx context.Context) (users []*types.User, err error)
	GetUser(ctx context.Context, userReq *types.User) (user *types.User, err error)
	UpdateUser(ctx context.Context, userReq *types.User) (user *types.User, err error)
	BlockUser(ctx context.Context, userId uint32) (ok bool, err error)
	UnblockUser(ctx context.Context, userId uint32) (ok bool, err error)

	CreateService(ctx context.Context, name string) (s *types.Service, err error)
	GetAllServices(ctx context.Context) (ss []*types.Service, err error)
	GetService(ctx context.Context, svc *types.Service) (s *types.Service, err error)

	CreateAccount(ctx context.Context) (a *types.Account, err error)
	CreateAccountWithName(ctx context.Context, name string) (a *types.Account, err error)
	GetAllAccounts(ctx context.Context) (as []*types.Account, err error)
	GetAccount(ctx context.Context, acc *types.Account) (a *types.Account, err error)
	UpdateAccount(ctx context.Context, acc *types.Account) (a *types.Account, err error)

	AttachUserToAccount(ctx context.Context, userId, accountId uint32) (ok bool, err error)
	AttachAccountToService(ctx context.Context, serviceId, accountId uint32) (ok bool, err error)
	RemoveAccountFromService(ctx context.Context, serviceId, accountId uint32) (ok bool, err error)

	CreatePermission(ctx context.Context, serviceId uint32, name string, access *access.Access) (p *types.Permission, err error)
	GetPermission(ctx context.Context, p *types.Permission) (perm *types.Permission, err error)
	GetAllPermission(ctx context.Context) (p []*types.Permission, err error)
	GetFilteredPermissions(ctx context.Context, p *types.Permission) (perm []*types.Permission, err error)
	DeletePermission(ctx context.Context, p *types.Permission) (ok bool, err error)

	GetUserPermissions(ctx context.Context, userId uint32) (permissions []*types.Permission, err error)
	AddUserPermission(ctx context.Context, p *types.Permission, userId uint32) (ok bool, err error)
	RemoveUserPermission(ctx context.Context, permId, userId uint32) (ok bool, err error)
}

type service struct {
	svcsrv.Service
	db         *gorm.DB
	l          *log.Logger
	user       user.Service
	permission permission.Service
	svc        svcsrv.Service
	acc        account.Repo
}

func New(l *log.Logger, db *gorm.DB, svc svcsrv.Service, acc account.Repo, p permission.Service, u user.Service) Service {
	return &service{l: l, db: db, svc: svc, acc: acc, permission: p, user: u}
}

func (s service) CreateUserWithLoginPassword(ctx context.Context, login, pass string) (*types.User, error) {
	return s.user.CreateWithLoginPassword(ctx, login, pass)
}

func (s service) CreateUserWithTelegram(ctx context.Context, id uint64, name, userN string) (*types.User, error) {
	return s.user.CreateWithTelegram(ctx, id, name, userN)
}

func (s service) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	return s.user.GetAll(ctx)
}

func (s service) GetUser(ctx context.Context, user *types.User) (*types.User, error) {
	return s.user.Get(ctx, user)
}

func (s service) UpdateUser(ctx context.Context, user *types.User) (*types.User, error) {
	return s.user.Update(ctx, user)
}

func (s service) UpdateMapUser(ctx context.Context, m map[string]interface{}) (bool, error) {
	return s.user.UpdateMap(ctx, m)
}

func (s service) BlockUser(ctx context.Context, userId uint32) (bool, error) {
	return s.user.Block(ctx, userId)
}

func (s service) UnblockUser(ctx context.Context, userId uint32) (bool, error) {
	return s.user.Unblock(ctx, userId)
}

func (s service) CreateService(ctx context.Context, name string) (*types.Service, error) {
	return s.svc.Create(ctx, name)
}

func (s service) GetAllServices(ctx context.Context) ([]*types.Service, error) {
	return s.svc.GetAll(ctx)
}

func (s service) GetService(ctx context.Context, svc *types.Service) (*types.Service, error) {
	return s.svc.Get(ctx, svc)
}

func (s service) CreateAccount(ctx context.Context) (*types.Account, error) {
	return s.acc.Create(ctx)
}

func (s service) CreateAccountWithName(ctx context.Context, name string) (*types.Account, error) {
	return s.acc.CreateWithName(ctx, name)
}

func (s service) GetAllAccounts(ctx context.Context) ([]*types.Account, error) {
	return s.acc.GetAll(ctx)
}

func (s service) GetAccount(ctx context.Context, acc *types.Account) (*types.Account, error) {
	return s.acc.Get(ctx, acc)
}

func (s service) UpdateAccount(ctx context.Context, acc *types.Account) (*types.Account, error) {
	return s.acc.Update(ctx, acc)
}

//func (s service) UpdateMapAccount(ctx context.Context, m map[string]interface{}) (bool, error) {
//	return s.acc.UpdateMap(ctx, m)
//}

func (s service) CreatePermission(ctx context.Context, serviceId uint32, name string, access *access.Access) (p *types.Permission, err error) {
	return s.permission.Create(ctx, serviceId, name, access)
}

func (s service) GetPermission(ctx context.Context, p *types.Permission) (perm *types.Permission, err error) {
	return s.permission.Get(ctx, p)
}

func (s service) GetAllPermission(ctx context.Context) (p []*types.Permission, err error) {
	return s.permission.GetAll(ctx)
}

func (s service) GetFilteredPermissions(ctx context.Context, p *types.Permission) (perm []*types.Permission, err error) {
	return s.permission.GetFiltered(ctx, p)
}

func (s service) DeletePermission(ctx context.Context, p *types.Permission) (ok bool, err error) {
	return s.permission.Delete(ctx, p)
}

func (s service) Create(ctx context.Context, name string) (*types.Service, error) {
	return s.Service.Create(ctx, name)
}

func (s service) GetAll(ctx context.Context) ([]*types.Service, error) {
	return s.Service.GetAll(ctx)
}

func (s service) Get(ctx context.Context, svc *types.Service) (*types.Service, error) {
	return s.Service.Get(ctx, svc)
}

func (s service) GetUserPermissions(ctx context.Context, userId uint32) (permissions []*types.Permission, err error) {
	err = s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: types.Model{ID: userId}}).
		Association("Permissions").
		Find(&permissions)
	s.l.Println("user permissions:")
	for i, p := range permissions {
		s.l.Println(i, p)
	}
	return
}

func (s service) AddUserPermission(ctx context.Context, p *types.Permission, userId uint32) (ok bool, err error) {
	err = s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: types.Model{ID: userId}}).
		Association("Permissions").
		Append(&types.Permission{Model: types.Model{ID: p.ID}})
	ok = err == nil
	return
}

func (s service) RemoveUserPermission(ctx context.Context, permId, userId uint32) (ok bool, err error) {
	err = s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: types.Model{ID: userId}}).
		Association("Permissions").
		Delete(&types.Permission{Model: types.Model{ID: permId}})
	ok = err == nil
	return
}

func (s service) AttachUserToAccount(ctx context.Context, userId, accountId uint32) (ok bool, err error) {
	u, err := s.user.Update(ctx, &types.User{Model: types.Model{ID: userId}, AccountID: accountId})
	if err != nil {
		return false, err
	}
	if u.ID != userId {
		return false, fmt.Errorf("updated userID=%d != userID=%d", u.ID, userId)
	}
	if u.AccountID != accountId {
		return false, fmt.Errorf("updated accountID=%d != accountID=%d", u.AccountID, accountId)
	}
	return true, nil
}

func (s service) AttachAccountToService(ctx context.Context, serviceId, accountID uint32) (bool, error) {
	s.l.Println("AttachAccountToService", serviceId, accountID)
	v := types.Account{Model: types.Model{ID: accountID}}
	svc := types.Service{Model: types.Model{ID: serviceId}}
	err := s.db.Debug().WithContext(ctx).
		Model(&svc).
		Association("Accounts").
		Append(&v)
	return err == nil, err
}

func (s service) RemoveAccountFromService(ctx context.Context, serviceId, accountId uint32) (bool, error) {
	s.l.Println("RemoveAccountFromService", serviceId, accountId)
	v := types.Account{Model: types.Model{ID: accountId}}
	svc := types.Service{Model: types.Model{ID: serviceId}}
	err := s.db.Debug().WithContext(ctx).
		Model(&svc).
		Association("Accounts").
		Delete(&v)
	return err == nil, err
}
