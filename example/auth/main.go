package main

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nullc4t/auth-service/pkg/access"
	"github.com/nullc4t/auth-service/pkg/auth"
	"github.com/nullc4t/auth-service/pkg/auth/proto"
	authgrpc "github.com/nullc4t/auth-service/pkg/auth/transport/grpc"
	transporthttp "github.com/nullc4t/auth-service/pkg/auth/transport/http"
	jwtservice "github.com/nullc4t/auth-service/pkg/jwt"
	"github.com/nullc4t/auth-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/url"
)

const (
	serviceAddr = "auth.dev.0x00.site"
	serviceName = "test"
	protocol    = "http"
)

func main() {
	/*
		init
	*/
	l := logger.New("[ example ]\t")

	var (
		client auth.Service
	)

	switch protocol {
	case "http":
		u, err := url.Parse("https://" + serviceAddr + "/")
		if err != nil {
			l.Fatal(err)
		}
		client = transporthttp.NewHTTPClient(u)
	case "grpc":
		conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			l.Fatal(err)
		}
		client = authgrpc.NewGRPCClient(conn, proto.Service_ServiceDesc.ServiceName)
	default:
		panic(protocol)
	}

	permissions, err := client.GetPermissionsForService(context.TODO(), serviceName)
	if err != nil {
		l.Fatal(err)
	}

	pNames := make([]string, 0, len(permissions))
	l.Println("service permissions")
	for _, permission := range permissions {
		l.Println(permission.ID, permission.ServiceID, permission.Name, permission.Access)
		pNames = append(pNames, permission.Name)
	}

	am := access.NewHelperFromPermissions(pNames...)
	l.Println("helper permissions")
	for name, acc := range am.ByName() {
		l.Println(name, acc)
	}

	l.Println("get public key")
	keyBytes, err := client.PublicKey(context.TODO())
	if err != nil {
		l.Fatal(err)
	}

	key, err := jwt.ParseECPublicKeyFromPEM(keyBytes)
	if err != nil {
		l.Fatal(err)
	}

	jwtSvc := jwtservice.New(
		logger.New("[ jwt service ]\t"),
		jwt.SigningMethodES256,
		jwtservice.ValidMethodsEC,
		jwtservice.AccessClaimsFactory,
		nil,
		key,
	)

	/*
		register
	*/
	ok, err := client.Register(
		context.TODO(),
		"login",
		"password",
		serviceName,
		0, // TODO
	)
	if err != nil {
		l.Fatal(err)
	}
	if !ok {
		l.Fatal("!ok")
	}

	/*
		login
	*/
	token, err := client.Login(context.TODO(), "login", "password", serviceName)
	if err != nil {
		l.Fatal(err)
	}
	l.Println(token.AccessToken)
	/*
		validate token
	*/
	cl, err := jwtSvc.VerifyAccessToken(token.AccessToken)
	if err != nil {
		l.Fatal(err)
	}

	/*
		validate access, must NOT validate
	*/
	claims := cl.(*jwtservice.AccessClaims)
	if jwtservice.ValidatorFactory(am, "test")(claims.Access) {
		l.Fatal("user has no permissions")
	}

	// go-kit middleware
	// authz middleware with test permissions required
	jwtservice.Middleware(jwtSvc, am, "test")
}
