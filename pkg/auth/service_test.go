package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	stdjwt "github.com/golang-jwt/jwt/v4"
	"github.com/nullc4t/auth-service/pkg/account"
	"github.com/nullc4t/auth-service/pkg/jwt"
	"github.com/nullc4t/auth-service/pkg/logger"
	"github.com/nullc4t/auth-service/pkg/mgmt"
	"github.com/nullc4t/auth-service/pkg/permission"
	svcsrv "github.com/nullc4t/auth-service/pkg/service"
	"github.com/nullc4t/auth-service/pkg/types"
	user3 "github.com/nullc4t/auth-service/pkg/user"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type testSuite struct {
	suite.Suite
	auth Service
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

	acco := account.NewLoggingMiddleware(
		logger.New("[ account auth ]\t"),
		account.New(db),
	)
	user := user3.NewLoggingMiddleware(logger.New("[ user auth ]\t"), user3.New(db))

	//authz := service2.LoggingMiddleware(
	//	log.NewLogfmtLogger(os.Stdout),
	//)(service2.RecoveringMiddleware(
	//	log.NewLogfmtLogger(os.Stdout),
	//)(authz2.NewService(db)))

	p := permission.New(db)

	svc := svcsrv.New(logger.New("[ service ]\t"), db)

	mgm := mgmt.New(logger.New("[ mgmt ]\t"), db, svc, acco, p, user)

	var privateKey = make([]byte, 64)
	_, err = rand.Read(privateKey)
	s.Require().NoError(err)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	s.Require().NoError(err)

	j := jwt.New(
		logger.New("[ jwt ]\t"),
		stdjwt.SigningMethodES256,
		[]string{stdjwt.SigningMethodES256.Name},
		jwt.AccessClaimsFactory,
		key,
		&key.PublicKey,
	)

	s.auth = New(logger.New("[ auth ]\t"), mgm, j)

	tables, err := db.Migrator().GetTables()
	s.Require().NoError(err)
	for _, table := range tables {
		s.Require().NoError(db.Debug().Migrator().DropTable(table))
	}
	s.Require().NoError(db.Debug().AutoMigrate(
		&types.User{}, &types.Account{}, &types.Service{}, &types.Permission{},
	))
}

func (s *testSuite) TestPrivateKey() {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	s.T().Log(pk)
	s.Require().NoError(err)
	b, err := x509.MarshalECPrivateKey(pk)
	//b, err := x509.MarshalPKCS8PrivateKey(pk)
	s.Require().NoError(err)
	s.T().Log(b)
	s.T().Log(string(b))

	pk, err = x509.ParseECPrivateKey(b)
	//pk, err = stdjwt.ParseECPrivateKeyFromPEM(b)
	s.Require().NoError(err)
	s.T().Log(pk)
}

func (s *testSuite) TestPublicKey() {
	pk, err := s.auth.PublicKey(context.TODO())
	s.NoError(err)
	s.T().Log(pk)
	s.T().Log(string(pk))
	pub, err := stdjwt.ParseECPublicKeyFromPEM(pk)
	s.Require().NoError(err)
	s.T().Log(pub)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(testSuite))
}
