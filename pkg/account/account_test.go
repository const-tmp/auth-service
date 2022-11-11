package account

import (
	"context"
	"fmt"
	"github.com/nullc4t/auth-service/pkg/logger"
	"github.com/nullc4t/auth-service/pkg/mgmt"
	"github.com/nullc4t/auth-service/pkg/types"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type testSuite struct {
	suite.Suite
	db   *gorm.DB
	acco Repo
	mgmt mgmt.Service
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
	s.acco = NewLoggingMiddleware(
		logger.New("[ account auth ]\t"),
		New(db),
	)
	s.mgmt = mgmt.New(logger.New("[ mgmt ]\t"), db)

	s.Require().NoError(s.db.Debug().Migrator().DropTable(&types.Account{}, &types.Service{}))
	s.Require().NoError(s.db.Debug().AutoMigrate(&types.Account{}, &types.Service{}))
}

func (s *testSuite) TestAccount() {
	s.Run("create", func() {
		res, err := s.acco.Create(context.TODO())
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("create name", func() {
		res, err := s.acco.CreateWithName(context.TODO(), "test")
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get all", func() {
		res, err := s.acco.GetAll(context.TODO())
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by name", func() {
		res, err := s.acco.Get(context.TODO(), &types.Account{Name: "test"})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by id", func() {
		res, err := s.acco.Get(context.TODO(), &types.Account{Model: types.Model{ID: 1}})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("update", func() {
		res, err := s.acco.Update(context.TODO(), &types.Account{Model: types.Model{ID: 2}, Name: "test test"})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by name", func() {
		res, err := s.acco.Get(context.TODO(), &types.Account{Name: "test test"})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get all", func() {
		res, err := s.acco.GetAll(context.TODO())
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by name", func() {
		res, err := s.acco.Get(context.TODO(), &types.Account{Name: "test"})
		s.T().Logf("%+v", res)
		s.T().Logf("%+v", err)
		s.Require().Error(err)
	})
}

func TestAccount(t *testing.T) {
	suite.Run(t, new(testSuite))
}
