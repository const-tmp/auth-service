// Code generated by microgen 1.0.5. DO NOT EDIT.

// It is better for you if you do not change functions names!
// This file will never be overwritten.
package transportgrpc

import (
	"github.com/gofrs/uuid"
	access "github.com/nullc4t/auth-service/pkg/access"
	pb "github.com/nullc4t/auth-service/pkg/mgmt/proto"
	types "github.com/nullc4t/auth-service/pkg/types"
)

func PtrTypesUserToProto(user *types.User) (*pb.User, error) {
	return &pb.User{
		Id:         user.ID,
		Code:       user.Code.String(),
		Name:       user.Name,
		TgName:     user.TGName,
		TgId:       user.TGID,
		TgUsername: user.TGUserName,
		ParentId:   user.ParentID,
		AccountId:  user.AccountID,
		Blocked:    user.Blocked,
	}, nil
}

func ProtoToPtrTypesUser(protoUser *pb.User) (*types.User, error) {
	u, err := uuid.FromString(protoUser.Code)
	if err != nil {
		return nil, err
	}
	return &types.User{
		Model: types.Model{
			ID: protoUser.Id,
		},
		Code:       u,
		Name:       protoUser.Name,
		TGName:     protoUser.TgName,
		TGID:       protoUser.TgId,
		TGUserName: protoUser.TgUsername,
		ParentID:   protoUser.ParentId,
		AccountID:  protoUser.AccountId,
		Blocked:    protoUser.Blocked,
	}, nil
}

func ListPtrTypesUserToProto(users []*types.User) ([]*pb.User, error) {
	res := make([]*pb.User, 0, len(users))
	for _, s := range users {
		res = append(res, &pb.User{
			Id:         s.ID,
			Code:       s.Code.String(),
			Name:       s.Name,
			TgName:     s.Name,
			TgId:       s.TGID,
			TgUsername: s.TGUserName,
			ParentId:   s.ParentID,
			AccountId:  s.AccountID,
			Blocked:    s.Blocked,
		})
	}
	return res, nil
}

func ProtoToListPtrTypesUser(protoUsers []*pb.User) ([]*types.User, error) {
	res := make([]*types.User, 0, len(protoUsers))
	for _, s := range protoUsers {
		u, err := uuid.FromString(s.Code)
		if err != nil {
			return nil, err
		}
		res = append(res, &types.User{
			Model:      types.Model{},
			Code:       u,
			Name:       s.Name,
			TGName:     s.TgName,
			TGID:       s.TgId,
			TGUserName: s.TgUsername,
			ParentID:   s.ParentId,
			AccountID:  s.AccountId,
			Blocked:    s.Blocked,
		})
	}
	return res, nil
}

func PtrTypesServiceToProto(s *types.Service) (*pb.Service, error) {
	return &pb.Service{
		Id:   s.ID,
		Name: s.Name,
		Code: s.Code.String(),
	}, nil
}

func ProtoToPtrTypesService(protoS *pb.Service) (*types.Service, error) {
	u, err := uuid.FromString(protoS.Code)
	if err != nil {
		return nil, err
	}
	return &types.Service{
		Model: types.Model{
			ID: protoS.Id,
		},
		Name: protoS.Name,
		Code: u,
	}, nil
}

func ListPtrTypesServiceToProto(ss []*types.Service) ([]*pb.Service, error) {
	res := make([]*pb.Service, 0, len(ss))
	for _, s := range ss {
		res = append(res, &pb.Service{
			Id:   s.ID,
			Name: s.Name,
			Code: s.Code.String(),
		})
	}
	return res, nil
}

func ProtoToListPtrTypesService(protoSs []*pb.Service) ([]*types.Service, error) {
	res := make([]*types.Service, 0, len(protoSs))
	for _, s := range protoSs {
		u, err := uuid.FromString(s.Code)
		if err != nil {
			return nil, err
		}
		res = append(res, &types.Service{
			Model: types.Model{
				ID: s.Id,
			},
			Name: s.Name,
			Code: u,
		})
	}
	return res, nil
}

func PtrTypesAccountToProto(a *types.Account) (*pb.Account, error) {
	return &pb.Account{
		Id:   a.ID,
		Name: a.Name,
		Code: a.Code.String(),
	}, nil
}

func ProtoToPtrTypesAccount(protoA *pb.Account) (*types.Account, error) {
	u, err := uuid.FromString(protoA.Code)
	if err != nil {
		return nil, err
	}
	return &types.Account{
		Model: types.Model{
			ID: protoA.Id,
		},
		Code: u,
		Name: protoA.Name,
	}, nil
}

func ListPtrTypesAccountToProto(as []*types.Account) ([]*pb.Account, error) {
	res := make([]*pb.Account, 0, len(as))
	for _, a := range as {
		res = append(res, &pb.Account{
			Id:   a.ID,
			Name: a.Name,
			Code: a.Code.String(),
		})
	}
	return res, nil
}

func ProtoToListPtrTypesAccount(protoAs []*pb.Account) ([]*types.Account, error) {
	res := make([]*types.Account, 0, len(protoAs))
	for _, a := range protoAs {
		u, err := uuid.FromString(a.Code)
		if err != nil {
			return nil, err
		}
		res = append(res, &types.Account{
			Model: types.Model{
				ID: a.Id,
			},
			Code: u,
			Name: a.Name,
		})
	}
	return res, nil
}

func PtrAccessAccessToProto(access *access.Access) (*pb.Access, error) {
	return &pb.Access{Access: uint64(*access)}, nil
}

func ProtoToPtrAccessAccess(protoAccess *pb.Access) (*access.Access, error) {
	a := access.Access(protoAccess.Access)
	return &a, nil
}

func PtrTypesPermissionToProto(p *types.Permission) (*pb.Permission, error) {
	return &pb.Permission{
		Id:        p.ID,
		Name:      p.Name,
		Access:    uint64(p.Access),
		ServiceId: p.ServiceID,
	}, nil
}

func ProtoToPtrTypesPermission(protoP *pb.Permission) (*types.Permission, error) {
	return &types.Permission{
		Model: types.Model{
			ID: protoP.Id,
		},
		ServiceID: protoP.ServiceId,
		Name:      protoP.Name,
		Access:    access.Access(protoP.Access),
	}, nil
}

func ListPtrTypesPermissionToProto(p []*types.Permission) ([]*pb.Permission, error) {
	res := make([]*pb.Permission, 0, len(p))
	for _, permission := range p {
		res = append(res, &pb.Permission{
			Id:        permission.ID,
			Name:      permission.Name,
			Access:    uint64(permission.Access),
			ServiceId: permission.ServiceID,
		})
	}
	return res, nil
}

func ProtoToListPtrTypesPermission(protoP []*pb.Permission) ([]*types.Permission, error) {
	res := make([]*types.Permission, 0, len(protoP))
	for _, permission := range protoP {
		res = append(res, &types.Permission{
			Model: types.Model{
				ID: permission.Id,
			},
			ServiceID: permission.ServiceId,
			Name:      permission.Name,
			Access:    access.Access(permission.Access),
		})
	}
	return res, nil
}
