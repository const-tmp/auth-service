package user

import (
	"auth/pkg/types"
	"context"
)

type Service interface {
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

type service struct {
	r Repo
}

func (s service) SetAccount(ctx context.Context, userID, accID uint) (bool, error) {
	return s.r.SetAccount(ctx, userID, accID)
}

func NewService(r Repo) Service {
	return &service{r: r}
}

func (s service) CreateWithLoginPassword(ctx context.Context, login, pass string) (types.User, error) {
	return s.r.CreateWithLoginPassword(ctx, login, pass)
}

func (s service) CreateWithTelegram(ctx context.Context, id uint64, name, userN string) (types.User, error) {
	return s.r.CreateWithTelegram(ctx, id, name, userN)
}

func (s service) GetAll(ctx context.Context) ([]types.User, error) {
	return s.r.GetAll(ctx)
}

func (s service) Get(ctx context.Context, user types.User) (types.User, error) {
	return s.r.Get(ctx, user)
}

func (s service) Update(ctx context.Context, user types.User) (types.User, error) {
	return s.r.Update(ctx, user)
}

func (s service) UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error) {
	return s.r.UpdateMap(ctx, m)
}

func (s service) Block(ctx context.Context, userID uint) (bool, error) {
	return s.r.Block(ctx, userID)
}

func (s service) Unblock(ctx context.Context, userID uint) (bool, error) {
	return s.r.Unblock(ctx, userID)
}
