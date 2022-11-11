package authz

//import (
//	"github.com/nullc4t/auth-service/pkg/types"
//	"context"
//	"gorm.io/gorm"
//)
//
//// @microgen middleware, logging, recovering
//type Service interface {
//	GetPermissions(ctx context.Context, serviceID uint32) (p []*types.Permission, err error)
//	AddPermission(ctx context.Context, perm *types.Permission) (p *types.Permission, err error)
//	RemovePermission(ctx context.Context, p *types.Permission) (ok bool, err error)
//	GetUserPermissions(ctx context.Context, p *types.Permission, userID uint32) (permissions []*types.Permission, err error)
//	AddUserPermission(ctx context.Context, permID, userID uint32) (ok bool, err error)
//	RemoveUserPermission(ctx context.Context, permID, userID uint32) (ok bool, err error)
//}
//
//type service struct {
//	db *gorm.DB
//}
//
//func NewService(db *gorm.DB) Service {
//	return service{db: db}
//}
//
//func (s service) GetPermissions(ctx context.Context, serviceID uint32) ([]*types.Permission, error) {
//	var res []*types.Permission
//	err := s.db.Debug().WithContext(ctx).
//		Model(&types.Service{Model: types.Model{ID: serviceID}}).
//		Association("Permissions").Find(&res)
//	return res, err
//}
//
//func (s service) AddPermission(ctx context.Context, p *types.Permission) (*types.Permission, error) {
//	err := s.db.Debug().WithContext(ctx).
//		Model(&types.Service{Model: types.Model{ID: p.ServiceID}}).
//		Association("Permissions").Append(&p)
//	return p, err
//}
//
//func (s service) RemovePermission(ctx context.Context, p *types.Permission) (bool, error) {
//	err := s.db.Debug().WithContext(ctx).
//		Model(&types.Service{Model: types.Model{ID: p.ServiceID}}).
//		Association("Permissions").Delete(&p)
//	return err == nil, err
//}
//
//func (s service) GetUserPermissions(ctx context.Context, p *types.Permission, userID uint32) ([]*types.Permission, error) {
//	var res []*types.Permission
//	err := s.db.Debug().WithContext(ctx).
//		Model(&types.User{Model: types.Model{ID: userID}}).
//		Where("service_id = ?", p.ServiceID).
//		Association("Permissions").Find(&res)
//	return res, err
//}
//
//func (s service) AddUserPermission(ctx context.Context, permID, userID uint32) (bool, error) {
//	err := s.db.Debug().WithContext(ctx).
//		Model(&types.User{Model: types.Model{ID: userID}}).
//		Association("Permissions").Append(&types.Permission{Model: types.Model{ID: permID}})
//	return err == nil, err
//}
//
//func (s service) RemoveUserPermission(ctx context.Context, permID uint32, userID uint32) (bool, error) {
//	err := s.db.Debug().WithContext(ctx).
//		Model(&types.User{Model: types.Model{ID: userID}}).
//		Association("Permissions").Delete(&permID)
//	return err == nil, err
//}
