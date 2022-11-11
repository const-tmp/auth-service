package service

import (
	"context"
	"github.com/nullc4t/auth-service/pkg/types"
	"gorm.io/gorm"
	"log"
)

type Service interface {
	Create(ctx context.Context, name string) (*types.Service, error)
	GetAll(ctx context.Context) ([]*types.Service, error)
	Get(ctx context.Context, svc *types.Service) (*types.Service, error)
}

type service struct {
	db *gorm.DB
	l  *log.Logger
}

func New(l *log.Logger, db *gorm.DB) Service {
	return &service{l: l, db: db}
}

func (s service) Create(ctx context.Context, name string) (*types.Service, error) {
	s.l.Println("Create", name)
	v := types.Service{Name: name}
	err := s.db.Debug().WithContext(ctx).Create(&v).Error
	return &v, err
}

func (s service) GetAll(ctx context.Context) ([]*types.Service, error) {
	s.l.Println("GetAll")
	var v []*types.Service
	err := s.db.Debug().WithContext(ctx).Find(&v).Error
	return v, err
}

func (s service) Get(ctx context.Context, svc *types.Service) (*types.Service, error) {
	s.l.Println("Get", svc)
	var v types.Service
	err := s.db.Debug().WithContext(ctx).
		Where(&svc).
		Preload("Permissions").
		First(&v).Error
	return &v, err
}
