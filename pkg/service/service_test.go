package service

import (
	"context"
	"fmt"
	"github.com/nullc4t/auth-service/pkg/types"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type testSuite struct {
	suite.Suite
	db  *gorm.DB
	svc Service
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
	s.svc = New(logger.New("[ account auth ]\t"), db)
	s.Require().NoError(s.db.Debug().Migrator().DropTable(&types.Service{}, &types.Permission{}))
	s.Require().NoError(s.db.Debug().AutoMigrate(&types.Service{}, &types.Permission{}))
}

func (s *testSuite) TestServices() {
	var svc types.Service
	s.Run("create", func() {
		sv, err := s.svc.Create(context.TODO(), "test")
		s.Require().NoError(err)
		svc = sv
	})
	s.Run("get", func() {
		v, err := s.svc.Get(context.TODO(), svc)
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
	})
	s.Run("get all", func() {
		v, err := s.svc.GetAll(context.TODO())
		s.Require().NoError(err)
		for i, t := range v {
			s.T().Logf("%d\t%+v", i, t)
		}
	})
}

func TestService(t *testing.T) {
	suite.Run(t, new(testSuite))
}
