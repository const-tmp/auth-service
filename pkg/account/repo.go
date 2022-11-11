package account

import (
	"context"
	"github.com/nullc4t/auth-service/pkg/types"
	"gorm.io/gorm"
)

type Repo interface {
	Create(ctx context.Context) (*types.Account, error)
	CreateWithName(ctx context.Context, name string) (*types.Account, error)
	GetAll(ctx context.Context) ([]*types.Account, error)
	Get(ctx context.Context, acc *types.Account) (*types.Account, error)
	Update(ctx context.Context, acc *types.Account) (*types.Account, error)
	//UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error)
}

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *repo {
	return &repo{db: db}
}

func (r repo) Create(ctx context.Context) (*types.Account, error) {
	a := types.Account{}
	err := r.db.Debug().WithContext(ctx).Create(&a).Error
	return &a, err
}

func (r repo) CreateWithName(ctx context.Context, name string) (*types.Account, error) {
	a := types.Account{Name: name}
	err := r.db.Debug().WithContext(ctx).Create(&a).Error
	return &a, err
}

func (r repo) GetAll(ctx context.Context) ([]*types.Account, error) {
	var a []*types.Account
	err := r.db.Debug().WithContext(ctx).Find(&a).Error
	return a, err
}

func (r repo) Get(ctx context.Context, acc *types.Account) (*types.Account, error) {
	err := r.db.Debug().WithContext(ctx).Where(&acc).First(&acc).Error
	return acc, err
}

func (r repo) Update(ctx context.Context, acc *types.Account) (*types.Account, error) {
	err := r.db.Debug().WithContext(ctx).Updates(&acc).Error
	return acc, err
}

func (r repo) UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error) {
	err := r.db.Debug().WithContext(ctx).Updates(m).Error
	return err == nil, err
}
