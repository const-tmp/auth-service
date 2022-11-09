package mgmt

import (
	"auth/pkg/types"
	"context"
	"gorm.io/gorm"
	"log"
)

type Service interface {
	Create(ctx context.Context, name string) (types.Service, error)
	GetServices(ctx context.Context) ([]types.Service, error)
	GetService(ctx context.Context, svc types.Service) (types.Service, error)
	AttachAccountToService(ctx context.Context, serviceID, accountID uint) (bool, error)
	RemoveAccountFromService(ctx context.Context, serviceID, accountID uint) (bool, error)
}

type service struct {
	db *gorm.DB
	l  *log.Logger
}

func New(l *log.Logger, db *gorm.DB) Service {
	return &service{l: l, db: db}
}

func (s service) Create(ctx context.Context, name string) (types.Service, error) {
	s.l.Println("Create", name)
	v := types.Service{Name: name}
	err := s.db.Debug().WithContext(ctx).Create(&v).Error
	return v, err
}

func (s service) GetServices(ctx context.Context) ([]types.Service, error) {
	s.l.Println("GetServices")
	var v []types.Service
	err := s.db.Debug().WithContext(ctx).Find(&v).Error
	return v, err
}

func (s service) GetService(ctx context.Context, svc types.Service) (types.Service, error) {
	s.l.Println("GetService", svc)
	var v types.Service
	err := s.db.Debug().WithContext(ctx).
		Where(&svc).
		Preload("Permissions").
		First(&v).Error
	return v, err
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
