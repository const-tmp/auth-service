package main

import (
	"context"
	"fmt"
	"github.com/nullc4t/auth-service/pkg/access"
	"github.com/nullc4t/auth-service/pkg/mgmt"
	"github.com/nullc4t/auth-service/pkg/mgmt/proto"
	transportgrpc "github.com/nullc4t/auth-service/pkg/mgmt/transport/grpc"
	"github.com/nullc4t/auth-service/pkg/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

const serviceAddr = "localhost:9091"

type testSuite struct {
	suite.Suite
	mgmtClient mgmt.Service
	db         *gorm.DB
}

func TestMgmt(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupSuite() {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		"localhost",
		"postgres",
		"password",
		"postgres",
		5432,
		"Europe/Kiev",
	)))
	s.Require().NoError(err)
	s.db = db

	tables, err := s.db.Migrator().GetTables()
	s.Require().NoError(err)
	for _, table := range tables {
		s.Require().NoError(s.db.Debug().Migrator().DropTable(table))
	}
	s.Require().NoError(s.db.Debug().AutoMigrate(
		&types.User{}, &types.Account{}, &types.Service{}, &types.Permission{},
	))

	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.mgmtClient = transportgrpc.NewGRPCClient(conn, proto.Mgmt_ServiceDesc.ServiceName)
}

func (s *testSuite) TestMgmt() {
	am := access.NewHelperFromPermissions("test")

	var (
		svc  *types.Service
		user *types.User
		acc  *types.Account
		ps   []*types.Permission
	)
	s.Run("create service", func() {
		v, err := s.mgmtClient.CreateService(context.TODO(), "test")
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		svc = v
	})
	s.Run("create user", func() {
		v, err := s.mgmtClient.CreateUserWithLoginPassword(context.TODO(), "test", "test")
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		user = v
	})
	s.Run("create account", func() {
		v, err := s.mgmtClient.CreateAccount(context.TODO())
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		acc = v
	})
	s.Run("attach user to account", func() {
		v, err := s.mgmtClient.AttachUserToAccount(context.TODO(), user.ID, acc.ID)
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		s.True(v)
	})
	s.Run("get user", func() {
		v, err := s.mgmtClient.GetUser(context.TODO(), &types.User{Model: types.Model{ID: user.ID}})
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		s.Equal(acc.ID, v.AccountID)
	})
	for name, access := range am.ByName() {
		s.Run("create permission", func() {
			v, err := s.mgmtClient.CreatePermission(context.TODO(), svc.ID, name, &access)
			s.Require().NoError(err)
			s.T().Logf("%+v", v)
			ps = append(ps, v)
		})
	}
	for _, p := range ps {
		s.Run("add user permission", func() {
			v, err := s.mgmtClient.AddUserPermission(context.TODO(), p, user.ID)
			s.Require().NoError(err)
			s.T().Logf("%+v", v)
			s.True(v)
		})
	}
	s.Run("get user permission", func() {
		v, err := s.mgmtClient.GetUserPermissions(context.TODO(), user.ID)
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		s.Len(v, len(ps))
	})
}
