package user

import (
	"auth/pkg/password"
	"auth/pkg/types"
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type Repo interface {
	CreateWithLoginPassword(ctx context.Context, login, pass string) (types.User, error)
	CreateWithTelegram(ctx context.Context, id uint64, name, userN string) (types.User, error)
	GetAll(ctx context.Context) ([]types.User, error)
	Get(ctx context.Context, user types.User) (types.User, error)
	Update(ctx context.Context, user types.User) (types.User, error)
	UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error)
	Block(ctx context.Context, userID uint) (bool, error)
	Unblock(ctx context.Context, userID uint) (bool, error)
	SetAccount(ctx context.Context, userID, accID uint) (bool, error)
}

type repo struct {
	DB *gorm.DB
	l  *log.Logger
}

func (r repo) SetAccount(ctx context.Context, userID, accID uint) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Update("account_id", accID).Error
	return err == nil, err
}

func New(db *gorm.DB, l *log.Logger) Repo {
	return &repo{DB: db, l: l}
}

func (r repo) CreateWithLoginPassword(ctx context.Context, login, pass string) (types.User, error) {
	ph, err := password.Hash(pass)
	if err != nil {
		return types.User{}, fmt.Errorf("password hash errpr :%w", err)
	}
	u := types.User{Name: login, Password: ph}
	stmt := r.DB.Debug().WithContext(ctx)
	err = stmt.Omit("TGID", "TGUserName", "AccountID").Create(&u).Error
	return u, err
}

func (r repo) CreateWithTelegram(ctx context.Context, id uint64, name, userN string) (types.User, error) {
	u := types.User{TGID: id, TGUserName: userN, TGName: name}
	stmt := r.DB.Debug().WithContext(ctx)
	err := stmt.Omit("AccountID", "Name", "Password").Create(&u).Error
	return u, err
}

func (r repo) Get(ctx context.Context, user types.User) (types.User, error) {
	err := r.DB.Debug().WithContext(ctx).Where(&user).Preload("Permissions").First(&user).Error
	return user, err
}

func (r repo) GetAll(ctx context.Context) ([]types.User, error) {
	var v []types.User
	err := r.DB.Debug().WithContext(ctx).Preload("Permissions").Find(&v).Error
	return v, err
}

func (r repo) Update(ctx context.Context, user types.User) (types.User, error) {
	err := r.DB.Debug().WithContext(ctx).Updates(&user).Error
	return user, err
}

func (r repo) UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).Updates(m).Error
	return err == nil, err
}

func (r repo) Block(ctx context.Context, userID uint) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Update("blocked", true).Error
	return err == nil, err
}

func (r repo) Unblock(ctx context.Context, userID uint) (bool, error) {
	err := r.DB.Debug().WithContext(ctx).
		Model(&types.User{Model: gorm.Model{ID: userID}}).
		Update("blocked", false).Error
	return err == nil, err
}
