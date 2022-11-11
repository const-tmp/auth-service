package authz

//
//type loggingMiddleware struct {
//	l    *log.Logger
//	next Service
//}
//
//func NewLoggingMiddleware(l *log.Logger, next Service) Service {
//	return &loggingMiddleware{l: l, next: next}
//}
//
//func (l loggingMiddleware) GetPermissions(ctx context.Context, serviceID uint32) ([]*types.Permission, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (l loggingMiddleware) AddPermission(ctx context.Context, p *types.Permission) (*types.Permission, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (l loggingMiddleware) RemovePermission(ctx context.Context, p *types.Permission) (bool, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (l loggingMiddleware) GetUserPermissions(ctx context.Context, p *types.Permission, userID uint32) ([]*types.Permission, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (l loggingMiddleware) AddUserPermission(ctx context.Context, permID, userID uint32) (bool, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (l loggingMiddleware) RemoveUserPermission(ctx context.Context, permID uint32, userID uint32) (bool, error) {
//	//TODO implement me
//	panic("implement me")
//}
