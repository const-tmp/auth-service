package user

import (
	"context"
	"github.com/nullc4t/auth-service/pkg/types"
)

type Service interface {
	CreateWithLoginPassword(ctx context.Context, login, pass string) (*types.User, error)
	CreateWithTelegram(ctx context.Context, id uint64, name, userN string) (*types.User, error)
	GetAll(ctx context.Context) ([]*types.User, error)
	Get(ctx context.Context, user *types.User) (*types.User, error)
	Update(ctx context.Context, user *types.User) (*types.User, error)
	UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error)
	Block(ctx context.Context, userID uint32) (bool, error)
	Unblock(ctx context.Context, userID uint32) (bool, error)
}
