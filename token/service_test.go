package jwt

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"infra/cfg"
	"infra/pkg/models"
	"infra/pkg/redis"
	"testing"
)

type jwtTestSuite struct {
	suite.Suite
	service *Service
}

func TestJwt(t *testing.T) {
	suite.Run(t, new(jwtTestSuite))
}

func (s *jwtTestSuite) SetupSuite() {
	err := cfg.ReadConfig()
	s.Require().NoError(err)
	viper.Set(cfg.RedisHost, "localhost")

	j, err := NewService(redis.NewAutoCfg())
	s.Require().NoError(err)
	s.service = j
}

func (s *jwtTestSuite) TestJWT() {
	var accessToken, refreshToken string
	var uid uint64 = 1123
	var username = "jwttest"

	s.Run("generate", func() {
		at, rt, err := s.service.NewSession(&models.User{
			ID:     &uid,
			Access: 1,
			Name:   username,
		})
		s.Require().NoError(err)
		s.T().Log("access", at)
		s.T().Log("refresh", rt)
		accessToken = at
		refreshToken = rt
	})

	s.Run("verify access", func() {
		ac, at, valid, err := s.service.VerifyAccessToken(accessToken)
		s.Require().NoError(err)
		s.Require().True(valid)
		s.T().Logf("%+v", ac)
		s.T().Logf("%+v", at)
	})

	s.Run("verify refresh", func() {
		rc, rt, valid, err := s.service.VerifyRefreshToken(refreshToken)
		s.Require().NoError(err)
		s.Require().True(valid)
		s.T().Logf("%+v", rc)
		s.T().Logf("%+v", rt)
	})

	s.Run("logout", func() {
		claims, _, isValid, err := s.service.VerifyRefreshToken(refreshToken)
		s.Require().NoError(err)
		s.Require().True(isValid)

		_, err = s.service.Logout(claims.ID)
		s.Require().NoError(err)

		_, _, isValid, err = s.service.VerifyAccessToken(accessToken)
		s.Require().Error(err)
		s.Require().False(isValid)

		_, _, isValid, err = s.service.VerifyRefreshToken(refreshToken)
		s.Require().Error(err)
		s.Require().False(isValid)
	})
}
