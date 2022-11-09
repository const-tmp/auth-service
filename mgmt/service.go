package mgmt

import (
	"auth/pkg/types"
	"context"
	"gorm.io/gorm"
	"log"
)

type Service interface {
	AttachAccountToService(ctx context.Context, serviceID, accountID uint) (bool, error)
	RemoveAccountFromService(ctx context.Context, serviceID, accountID uint) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint) (permissions []types.Permission, err error)
	AddUserPermission(ctx context.Context, p types.Permission, userID uint) (ok bool, err error)
	RemoveUserPermission(ctx context.Context, permID, userID uint) (ok bool, err error)
}

type service struct {
	db *gorm.DB
	l  *log.Logger
}

func New(l *log.Logger, db *gorm.DB) Service {
	return &service{l: l, db: db}
}

func (s service) GetUserPermissions(ctx context.Context, userID uint) (permissions []types.Permission, err error) {
	err = s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Association("Permissions").
		Find(&permissions)
	return
}

func (s service) AddUserPermission(ctx context.Context, p types.Permission, userID uint) (ok bool, err error) {
	err = s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Association("Permissions").
		Append(&p)
	ok = err == nil
	return
}

func (s service) RemoveUserPermission(ctx context.Context, permID, userID uint) (ok bool, err error) {
	err = s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Association("Permissions").
		Delete(&types.Permission{Model: gorm.Model{ID: permID}})
	ok = err == nil
	return
}

func (s service) AttachAccountToService(ctx context.Context, serviceID uint, accountID uint) (bool, error) {
	s.l.Println("AttachAccountToService", serviceID, accountID)
	v := types.Account{Model: gorm.Model{ID: accountID}}
	svc := types.Service{Model: gorm.Model{ID: serviceID}}
	err := s.db.Debug().WithContext(ctx).
		Model(&svc).
		Association("Accounts").
		Append(&v)
	return err == nil, err
}

func (s service) RemoveAccountFromService(ctx context.Context, serviceID uint, accountID uint) (bool, error) {
	s.l.Println("RemoveAccountFromService", serviceID, accountID)
	v := types.Account{Model: gorm.Model{ID: accountID}}
	svc := types.Service{Model: gorm.Model{ID: serviceID}}
	err := s.db.Debug().WithContext(ctx).
		Model(&svc).
		Association("Accounts").
		Delete(&v)
	return err == nil, err
}
