package permission

import (
	"auth/access"
	"auth/logger"
	"auth/pkg/types"
	"context"

	"gorm.io/gorm"
	"log"
)

type Service interface {
	Create(ctx context.Context, serviceID uint32, name string, access *access.Access) (p *types.Permission, err error)
	Get(ctx context.Context, p *types.Permission) (perm *types.Permission, err error)
	GetAll(ctx context.Context) (p []*types.Permission, err error)
	GetFiltered(ctx context.Context, p *types.Permission) (perm []*types.Permission, err error)
	Delete(ctx context.Context, p *types.Permission) (ok bool, err error)
}

type service struct {
	db *gorm.DB
	l  *log.Logger
}

func New(db *gorm.DB) Service {
	return &service{l: logger.New("[ permission ] "), db: db}
}

func (s service) Create(ctx context.Context, serviceID uint32, name string, access *access.Access) (p *types.Permission, err error) {
	p = &types.Permission{
		ServiceID: serviceID,
		Name:      name,
		Access:    *access,
	}
	err = s.db.Debug().WithContext(ctx).Create(&p).Error
	return
}

func (s service) Get(ctx context.Context, p *types.Permission) (perm *types.Permission, err error) {
	err = s.db.Debug().WithContext(ctx).Where(&p).First(&perm).Error
	return
}

func (s service) GetAll(ctx context.Context) (p []*types.Permission, err error) {
	err = s.db.Debug().WithContext(ctx).Find(&p).Error
	return
}

func (s service) GetFiltered(ctx context.Context, p *types.Permission) (perm []*types.Permission, err error) {
	err = s.db.Debug().WithContext(ctx).Where(&p).Find(&perm).Error
	return
}

func (s service) Delete(ctx context.Context, p *types.Permission) (ok bool, err error) {
	err = s.db.Debug().WithContext(ctx).Where(&p).Delete(&p).Error
	return
}
