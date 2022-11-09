package authz

import (
	"auth/pkg/types"
	"context"
	"gorm.io/gorm"
)

// @microgen middleware, logging, recovering
type Service interface {
	GetPermissions(ctx context.Context, serviceID uint) (p []types.Permission, err error)
	AddPermission(ctx context.Context, perm types.Permission) (p types.Permission, err error)
	RemovePermission(ctx context.Context, p types.Permission) (ok bool, err error)
	GetUserPermissions(ctx context.Context, p types.Permission, userID uint) (permissions []types.Permission, err error)
	AddUserPermission(ctx context.Context, permID, userID uint) (ok bool, err error)
	RemoveUserPermission(ctx context.Context, permID, userID uint) (ok bool, err error)
}

type service struct {
	db *gorm.DB
}

func New(db *gorm.DB) Service {
	return service{db: db}
}

func (s service) GetPermissions(ctx context.Context, serviceID uint) ([]types.Permission, error) {
	var res []types.Permission
	err := s.db.Debug().WithContext(ctx).
		Model(&types.Service{Model: gorm.Model{ID: serviceID}}).
		Association("Permissions").Find(&res)
	return res, err
}

func (s service) AddPermission(ctx context.Context, p types.Permission) (types.Permission, error) {
	err := s.db.Debug().WithContext(ctx).
		Model(&types.Service{Model: gorm.Model{ID: p.ServiceID}}).
		Association("Permissions").Append(&p)
	return p, err
}

func (s service) RemovePermission(ctx context.Context, p types.Permission) (bool, error) {
	err := s.db.Debug().WithContext(ctx).
		Model(&types.Service{Model: gorm.Model{ID: p.ServiceID}}).
		Association("Permissions").Delete(&p)
	return err == nil, err
}

func (s service) GetUserPermissions(ctx context.Context, p types.Permission, userID uint) ([]types.Permission, error) {
	var res []types.Permission
	err := s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Where("service_id = ?", p.ServiceID).
		Association("Permissions").Find(&res)
	return res, err
}

func (s service) AddUserPermission(ctx context.Context, permID, userID uint) (bool, error) {
	err := s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Association("Permissions").Append(&types.Permission{Model: gorm.Model{ID: permID}})
	return err == nil, err
}

func (s service) RemoveUserPermission(ctx context.Context, permID uint, userID uint) (bool, error) {
	err := s.db.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Association("Permissions").Delete(&permID)
	return err == nil, err
}
