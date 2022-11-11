package mgmt

import (
	"auth/account"
	"auth/logger"
	"auth/permission"
	svcsrv "auth/pkg/service"
	"auth/pkg/types"
	user2 "auth/pkg/user"
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type testSuite struct {
	suite.Suite
	db   *gorm.DB
	mgmt Service
	acc  account.Repo
	perm permission.Service
	user user2.Service
	svc  svcsrv.Service
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
	s.mgmt = New(logger.New("[ mgmt ]\t"), db)
	s.user = user2.NewLoggingMiddleware(logger.New("[ user ]\t"), user2.New(db))

	s.acc = account.NewLoggingMiddleware(
		logger.New("[ account ]\t"),
		account.New(db),
	)
	s.perm = permission.New(db)
	s.svc = svcsrv.New(
		logger.New("[ account ]\t"),
		db,
	)
	//s.Require().NoError(s.db.Debug().Migrator().DropTable(&types.Service{}, &types.Permission{}, &types.User{}, &types.Account{}))
	tables, err := s.db.Migrator().GetTables()
	s.Require().NoError(err)
	for _, table := range tables {
		s.Require().NoError(s.db.Debug().Migrator().DropTable(table))
	}
	s.Require().NoError(s.db.Debug().AutoMigrate(&types.Service{}, &types.Permission{}, &types.User{}, &types.Account{}))
}

func (s *testSuite) TestAccountService() {
	svc, err := s.svc.Create(context.TODO(), "test auth")
	s.Require().NoError(err)

	acc, err := s.acc.Create(context.TODO())
	s.Require().NoError(err)

	ok, err := s.mgmt.AttachAccountToService(context.TODO(), svc.ID, acc.ID)
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s *testSuite) TestPermission() {
	svc, err := s.svc.Create(context.TODO(), "test auth 2")
	s.Require().NoError(err)

	per, err := s.perm.Create(context.TODO(), svc.ID, "test", 1)
	s.Require().NoError(err)
	s.T().Logf("%+v", per)

	u, err := s.user.CreateWithTelegram(context.TODO(), 1, "test", "test")
	s.Require().NoError(err)
	s.T().Logf("%+v", u)

	ok, err := s.mgmt.AddUserPermission(context.TODO(), per, u.ID)
	s.Require().NoError(err)
	s.Require().True(ok)

	u, err = s.user.Get(context.TODO(), types.User{Model: types.Model{ID: u.ID}})
	s.Require().NoError(err)
	s.T().Logf("%+v", u)
	s.T().Logf("%+v", u.Permissions)
}

func TestService(t *testing.T) {
	suite.Run(t, new(testSuite))
}
