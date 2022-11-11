package permission

import (
	"auth/logger"
	svcsrv "auth/pkg/service"
	"auth/pkg/types"
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
	perm Service
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
	s.perm = New(db)
	s.svc = svcsrv.New(logger.New("[ auth ]\tt"), db)
	s.Require().NoError(s.db.Debug().Migrator().DropTable(&types.Permission{}, &types.Service{}))
	s.Require().NoError(s.db.Debug().AutoMigrate(&types.Permission{}, &types.Service{}))
}

func (s *testSuite) TestServices() {
	var svc types.Service

	s.Run("create auth", func() {
		v, err := s.svc.Create(context.TODO(), "test")
		s.Require().NoError(err)
		s.T().Logf("%+v", v)
		svc = v
	})

	testCases := []*types.Permission{
		{
			ServiceID: svc.ID,
			Name:      "test",
			Access:    1,
		},
	}
	for _, testCase := range testCases {
		s.Run("create", func() {
			v, err := s.perm.Create(context.TODO(), testCase.ServiceID, testCase.Name, testCase.Access)
			s.Require().NoError(err)
			s.T().Logf("%+v", v)
		})
		s.Run("get", func() {
			v, err := s.perm.Get(context.TODO(), testCase)
			s.Require().NoError(err)
			s.T().Logf("%+v", v)
		})
		s.Run("get all", func() {
			v, err := s.perm.GetAll(context.TODO())
			s.Require().NoError(err)
			for i, t := range v {
				s.T().Logf("%d\t%+v", i, t)
			}
		})
	}
}

func TestService(t *testing.T) {
	suite.Run(t, new(testSuite))
}
