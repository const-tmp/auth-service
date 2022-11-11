package user

import (
	"auth/pkg/types"
	"context"
	"log"
)

type loggingMiddleware struct {
	l    *log.Logger
	next Service
}

func NewLoggingMiddleware(l *log.Logger, next Service) Service {
	return &loggingMiddleware{l: l, next: next}
}

func (l loggingMiddleware) Block(ctx context.Context, userID uint32) (bool, error) {
	l.l.Println("Method: Block\tArgs:", userID)
	res, err := l.next.Block(ctx, userID)
	if err != nil {
		l.l.Println("Method: Block\tError:", err)
	} else {
		l.l.Println("Method: Block\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) Unblock(ctx context.Context, userID uint32) (bool, error) {
	l.l.Println("Method: Unblock\tArgs:", userID)
	res, err := l.next.Unblock(ctx, userID)
	if err != nil {
		l.l.Println("Method: Unblock\tError:", err)
	} else {
		l.l.Println("Method: Unblock\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) CreateWithLoginPassword(ctx context.Context, login, pass string) (*types.User, error) {
	l.l.Println("Method: CreateWithLoginPassword\tArgs:", login, pass)
	res, err := l.next.CreateWithLoginPassword(ctx, login, pass)
	if err != nil {
		l.l.Println("Method: CreateWithLoginPassword\tError:", err)
	} else {
		l.l.Println("Method: CreateWithLoginPassword\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) CreateWithTelegram(ctx context.Context, id uint64, name, userN string) (*types.User, error) {
	l.l.Println("Method: CreateWithTelegram\tArgs:", id, name, userN)
	res, err := l.next.CreateWithTelegram(ctx, id, name, userN)
	if err != nil {
		l.l.Println("Method: CreateWithTelegram\tError:", err)
	} else {
		l.l.Println("Method: CreateWithTelegram\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) GetAll(ctx context.Context) ([]*types.User, error) {
	l.l.Println("Method: GetAll\tArgs:")
	res, err := l.next.GetAll(ctx)
	if err != nil {
		l.l.Println("Method: GetAll\tError:", err)
	} else {
		l.l.Println("Method: GetAll\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) Get(ctx context.Context, user *types.User) (*types.User, error) {
	l.l.Println("Method: Get\tArgs:", user)
	res, err := l.next.Get(ctx, user)
	if err != nil {
		l.l.Println("Method: Get\tError:", err)
	} else {
		l.l.Println("Method: Get\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) Update(ctx context.Context, user *types.User) (*types.User, error) {
	l.l.Println("Method: Update\tArgs:", user)
	res, err := l.next.Update(ctx, user)
	if err != nil {
		l.l.Println("Method: Update\tError:", err)
	} else {
		l.l.Println("Method: Update\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error) {
	l.l.Println("Method: UpdateMap\tArgs:", m)
	res, err := l.next.UpdateMap(ctx, m)
	if err != nil {
		l.l.Println("Method: UpdateMap\tError:", err)
	} else {
		l.l.Println("Method: UpdateMap\tResult:", res)
	}
	return res, err
}
