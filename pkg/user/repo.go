package user

import (
	"context"
	"fmt"
	"github.com/nullc4t/auth-service/pkg/password"
	"github.com/nullc4t/auth-service/pkg/types"
	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Service {
	return &repo{DB: db}
}

func (r repo) CreateWithLoginPassword(ctx context.Context, login, pass string) (*types.User, error) {
	ph, err := password.Hash(pass)
	if err != nil {
		return nil, fmt.Errorf("password hash errpr :%w", err)
	}
	u := types.User{Name: login, Password: ph}
	stmt := r.DB.Debug().WithContext(ctx)
	err = stmt.Omit("TGID", "TGUserName", "AccountID").Create(&u).Error
	return &u, err
}

func (r repo) CreateWithTelegram(ctx context.Context, id uint64, name, userN string) (*types.User, error) {
	u := types.User{TGID: id, TGUserName: userN, TGName: name}
	stmt := r.DB.Debug().WithContext(ctx)
	err := stmt.Omit("AccountID", "Name", "Password").Create(&u).Error
	return &u, err
}

func (r repo) Get(ctx context.Context, user *types.User) (*types.User, error) {
	err := r.DB.Debug().WithContext(ctx).Where(&user).Preload("Permissions").First(&user).Error
	return user, err
}

func (r repo) GetAll(ctx context.Context) ([]*types.User, error) {
	var v []*types.User
	err := r.DB.Debug().WithContext(ctx).Preload("Permissions").Find(&v).Error
	return v, err
}

func (r repo) Update(ctx context.Context, user *types.User) (*types.User, error) {
	err := r.DB.Debug().WithContext(ctx).Updates(&user).Error
	return user, err
}

func (r repo) UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).Updates(m).Error
	return err == nil, err
}

func (r repo) Block(ctx context.Context, userID uint32) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).
		Model(&types.User{Model: types.Model{ID: userID}}).
		Update("blocked", true).Error
	return err == nil, err
}

func (r repo) Unblock(ctx context.Context, userID uint32) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).
		Model(&types.User{Model: types.Model{ID: userID}}).
		Update("blocked", false).Error
	return err == nil, err
}
