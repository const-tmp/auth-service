package account

import (
	"auth/pkg/types"
	"context"
	"log"
)

type loggingMiddleware struct {
	l    *log.Logger
	next Repo
}

func NewLoggingMiddleware(l *log.Logger, next Repo) Repo {
	return &loggingMiddleware{l: l, next: next}
}

func (l loggingMiddleware) Create(ctx context.Context) (*types.Account, error) {
	l.l.Println("Method: Create\tArgs:")
	res, err := l.next.Create(ctx)
	if err != nil {
		l.l.Println("Method: Create\tError:", err)
	} else {
		l.l.Println("Method: Create\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) CreateWithName(ctx context.Context, name string) (*types.Account, error) {
	l.l.Println("Method: CreateWithName\tArgs:", name)
	res, err := l.next.CreateWithName(ctx, name)
	if err != nil {
		l.l.Println("Method: CreateWithName\tError:", err)
	} else {
		l.l.Println("Method: CreateWithName\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) GetAll(ctx context.Context) ([]*types.Account, error) {
	l.l.Println("Method: GetAll\tArgs:")
	res, err := l.next.GetAll(ctx)
	if err != nil {
		l.l.Println("Method: GetAll\tError:", err)
	} else {
		l.l.Println("Method: GetAll\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) Get(ctx context.Context, acc *types.Account) (*types.Account, error) {
	l.l.Println("Method: Get\tArgs:", acc)
	res, err := l.next.Get(ctx, acc)
	if err != nil {
		l.l.Println("Method: Get\tError:", err)
	} else {
		l.l.Println("Method: Get\tResult:", res)
	}
	return res, err
}

func (l loggingMiddleware) Update(ctx context.Context, acc *types.Account) (*types.Account, error) {
	l.l.Println("Method: Update\tArgs:", acc)
	res, err := l.next.Update(ctx, acc)
	if err != nil {
		l.l.Println("Method: Update\tError:", err)
	} else {
		l.l.Println("Method: Update\tResult:", res)
	}
	return res, err
}

//func (l loggingMiddleware) UpdateMap(ctx context.Context, m map[string]interface{}) (bool, error) {
//	l.l.Println("Method: UpdateMap\tArgs:", m)
//	res, err := l.next.UpdateMap(ctx, m)
//	if err != nil {
//		l.l.Println("Method: UpdateMap\tError:", err)
//	} else {
//		l.l.Println("Method: UpdateMap\tResult:", res)
//	}
//	return res, err
//}
