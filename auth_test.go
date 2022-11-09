package auth

import (
	"auth/account"
	"auth/auth"
	authz2 "auth/auth/authz"
	"auth/auth/authz/service"
	"auth/jwt"
	"auth/logger"
	"auth/mgmt"
	"auth/pkg/types"
	user2 "auth/user"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/nullc4ts/bitmask_authz/authz"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

type testSuite struct {
	suite.Suite
	db    *gorm.DB
	acco  account.Repo
	user  user2.Service
	authz authz2.Service
	mgmt  mgmt.Service
	auth  auth.Service
	jwt   jwt.Service
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
	s.acco = account.NewLoggingMiddleware(
		logger.New("[ account service ]\t"),
		account.New(db),
	)
	s.user = user2.NewLoggingMiddleware(
		logger.New("[ user service ]\t"),
		user2.NewService(user2.New(db, logger.New("[ user service ]\t"))),
	)
	s.authz = service.LoggingMiddleware(
		log.NewLogfmtLogger(os.Stdout),
	)(service.RecoveringMiddleware(
		log.NewLogfmtLogger(os.Stdout),
	)(authz2.New(db)),
	)
	s.mgmt = mgmt.New(logger.New("[ mgmt ]\t"), db)

	s.Require().NoError(s.db.Debug().Migrator().DropTable(
		&types.User{}, &types.Account{}, &types.Service{}, &types.Permission{},
	))
	s.Require().NoError(s.db.Debug().AutoMigrate(
		&types.User{}, &types.Account{}, &types.Service{}, &types.Permission{},
	))

	var privateKey = make([]byte, 64)
	_, err = rand.Read(privateKey)
	s.Require().NoError(err)

	//key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//s.Require().NoError(err)

	s.jwt = jwt.NewService(time.Minute*5, time.Hour*24, privateKey)
	s.auth = auth.New(logger.New("[ auth ]\t"), s.user, s.authz, s.mgmt, s.acco, s.jwt)
}

func (s *testSuite) TestAuth() {
	az := authz.New("read", "write")
	var svc types.Service
	s.Run("create service", func() {
		sv, err := s.mgmt.Create(context.TODO(), "test")
		s.Require().NoError(err)
		svc = sv
	})
	var perms []types.Permission
	for name, access := range az.ByName() {
		s.Run(fmt.Sprint("create permission", name, access, "service:", svc.Name, svc.ID), func() {
			v, err := s.authz.AddPermission(context.TODO(), types.Permission{
				ServiceID: svc.ID,
				Name:      name,
				Access:    access,
			})
			s.Require().NoError(err)
			s.T().Log(v.ID, v.Name, v.Access, v.ServiceID)
			perms = append(perms, v)
		})
	}
	svcNames := []string{"", svc.Name}
	accIDs := []uint{0, 1}
	var j int
	for i := 0; i < 2; i++ {
		for _, name := range svcNames {
			for _, accID := range accIDs {
				j++
				s.Run(fmt.Sprint("register", i, name, accID), func() {
					ok, err := s.auth.Register(context.TODO(), fmt.Sprint("test", j), "test", name, accID)
					s.Require().NoError(err)
					s.Require().True(ok)
				})

				for _, perm := range perms {
					s.Run("set user permissions", func() {
						u, err := s.user.Get(context.TODO(), types.User{Name: fmt.Sprint("test", j)})
						s.Require().NoError(err)
						s.T().Log(u.ID, u.AccountID, u.Name, u.Permissions)
						v, err := s.authz.AddUserPermission(context.TODO(), perm.ID, u.ID)
						s.Require().NoError(err)
						s.Require().True(v)
					})
				}

				s.Run(fmt.Sprint("login", i, name, accID), func() {
					t, err := s.auth.Login(context.TODO(), fmt.Sprint("test", j), "test", name)
					s.Require().NoError(err)
					s.T().Log(t.AccessToken)
					s.T().Log(t.RefreshToken)
					claims, token, valid, err := s.jwt.VerifyAccessToken(t.AccessToken)
					s.Require().NoError(err)
					s.Require().True(valid)
					s.T().Log(claims)
					s.T().Log(token)
				})
			}
		}
	}
}
func (s *testSuite) TestAccount() {
	//s.Run("services", func() {
	//	var svcs []types.Service
	//	s.Run("create", func() {
	//		svc := types.Service{
	//			Name: fmt.Sprintf("service 1"),
	//		}
	//		s.Require().NoError(s.db.Debug().Create(&svc).Error)
	//		s.T().Logf("%+v", svc)
	//		svcs = append(svcs, svc)
	//	})
	//
	//	s.Run("permissions", func() {
	//		for _, service := range svcs {
	//			for n, a := range authz.New("active", "read", "write").ByName() {
	//				s.Require().NoError(s.db.Debug().Create(&types.Permission{
	//					ServiceID: service.ID,
	//					Name:      n,
	//					Access:    a,
	//				}).Error)
	//			}
	//		}
	//	})
	//})
	//
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
		res, err := s.acco.Get(context.TODO(), types.Account{Name: "test"})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by id", func() {
		res, err := s.acco.Get(context.TODO(), types.Account{Model: gorm.Model{ID: 1}})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("update", func() {
		res, err := s.acco.Update(context.TODO(), types.Account{Model: gorm.Model{ID: 2}, Name: "test test"})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by name", func() {
		res, err := s.acco.Get(context.TODO(), types.Account{Name: "test test"})
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get all", func() {
		res, err := s.acco.GetAll(context.TODO())
		s.Require().NoError(err)
		s.T().Logf("%+v", res)
	})
	s.Run("get by name", func() {
		res, err := s.acco.Get(context.TODO(), types.Account{Name: "test"})
		s.T().Logf("%+v", res)
		s.T().Logf("%+v", err)
		s.Require().Error(err)
	})

	/*	s.Run("get service", func() {
			svc := types.Service{}
			s.Require().NoError(s.db.Debug().
				Preload("Permissions").
				Preload("Accounts").
				First(&svc, 1).Error)
			s.T().Logf("%+v", svc)
		})
		s.Run("get service", func() {
			svc := types.Service{}
			s.Require().NoError(s.db.Debug().
				Preload("Permissions").
				Preload("Accounts").
				First(&svc, 1).Error)
			s.T().Logf("%+v", svc)
		})
	*/
	var testUser types.User

	s.Run("create user", func() {
		u, err := s.user.CreateWithLoginPassword(context.TODO(), "test", "test")
		s.Require().NoError(err)
		testUser = u
		s.T().Logf("%+v", testUser)
	})
}

func (s *testSuite) TestService() {
	az := authz.New("active", "read", "write", "admin", "root")
	var testServices []types.Service
	for i := 0; i < 5; i++ {
		s.Run(fmt.Sprint("create service", i+1), func() {
			v, err := s.mgmt.Create(context.TODO(), fmt.Sprint("service", i+1))
			s.Require().NoError(err)
			s.T().Logf("%+v", v)
			testServices = append(testServices, v)
		})
	}

	for i, testService := range testServices {
		for name, access := range az.ByName() {
			s.Run(fmt.Sprint(i, testService.Name, "add permission", name), func() {
				v, err := s.authz.AddPermission(context.TODO(), types.Permission{
					ServiceID: testService.ID,
					Name:      name,
					Access:    access,
				})
				s.Require().NoError(err)
				s.T().Logf("%d\t%+v", i, v)
			})
		}
	}

	var testPermissions []types.Permission

	for i, svc := range testServices {
		s.Run(fmt.Sprint(i, "get", svc.Name, "permissions"), func() {
			v, err := s.authz.GetPermissions(context.TODO(), svc.ID)
			s.Require().NoError(err)
			s.T().Logf("%d\t%+v", i, v)
			testPermissions = append(testPermissions, v...)
		})
	}

	var testUsers []types.User

	for i := 0; i < 5; i++ {
		s.Run(fmt.Sprint("create user", i+1), func() {
			name := fmt.Sprint("user", i)
			v, err := s.user.CreateWithLoginPassword(context.TODO(), name, name)
			s.Require().NoError(err)
			s.T().Logf("%d\t%+v", i, v)
			testUsers = append(testUsers, v)
		})
	}

	for _, user := range testUsers {
		for _, svc := range testServices {
			p, err := s.authz.GetPermissions(context.TODO(), svc.ID)
			s.Require().NoError(err)
			for _, permission := range p {
				s.Run(fmt.Sprint("add permission", user.ID, svc.ID, permission.Name), func() {
					v, err := s.authz.AddUserPermission(context.TODO(), permission.ID, user.ID)
					s.Require().NoError(err)
					s.T().Logf("%+v", v)
				})
			}
		}
	}

	for name, access := range az.ByName() {
		s.T().Log(name, access)
	}

	for _, user := range testUsers {
		for _, svc := range testServices {
			s.Run(fmt.Sprint("get", user.Name, svc.Name, "permissions"), func() {
				v, err := s.authz.GetUserPermissions(context.TODO(), types.Permission{ServiceID: svc.ID}, user.ID)
				s.Require().NoError(err)
				s.T().Log("user", user.Name, user.ID)
				for i, permission := range v {
					s.T().Log(i, permission.ID, permission.ServiceID, permission.Name, permission.Access)
				}
			})
		}
	}

	testUsers, err := s.user.GetAll(context.TODO())
	s.Require().NoError(err)
	for i, user := range testUsers {
		s.T().Log(i, user.ID, user.Name, user.Permissions)
		for _, permission := range user.Permissions {
			s.T().Log(permission.ID, permission.ServiceID, permission.Name, permission.Access)
		}
	}
	//for _, user := range testUsers {
	//	s.Run("delete permission", func() {
	//		ok, err := s.authz.RemoveUserPermission(context.TODO(), types.Permission{Name: "root"}, user.ID)
	//		s.Require().NoError(err)
	//		s.Require().True(ok)
	//	})
	//}
	testUsers, err = s.user.GetAll(context.TODO())
	s.Require().NoError(err)
	for i, user := range testUsers {
		s.T().Log(i, user.ID, user.Name, user.Permissions)
		for _, permission := range user.Permissions {
			s.T().Log(permission.ID, permission.ServiceID, permission.Name, permission.Access)
		}
	}
}

func (s *testSuite) TestPermissions() {
	az := authz.New("active", "read", "write", "admin", "root")

	svc, err := s.mgmt.Create(context.TODO(), "svc1")
	s.Require().NoError(err)

	svc, err = s.mgmt.GetService(context.TODO(), svc)
	s.Require().NoError(err)
	s.T().Log(svc)

	p, err := s.authz.AddPermission(
		context.TODO(),
		types.Permission{ServiceID: svc.ID, Name: "active", Access: az.ByName()["active"]},
	)
	s.Require().NoError(err)
	s.T().Log(p)

	ps, err := s.authz.GetPermissions(context.TODO(), svc.ID)
	s.Require().NoError(err)
	s.T().Log(ps)

}

func TestAuth(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) TestInterface() {
	u := user2.New(s.db, logger.New(""))
	var _ user2.Repo = u
	var _ user2.Service = u
}
