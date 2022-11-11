// Microgen appends missed functions.
package main

import (
	"auth/account"
	"auth/logger"
	"auth/mgmt"
	"auth/mgmt/proto"
	mgmtservice "auth/mgmt/service"
	mgmttransport "auth/mgmt/transport"
	mgmtgrpc "auth/mgmt/transport/grpc"
	mgmthttp "auth/mgmt/transport/http"
	"auth/permission"
	"auth/pkg/types"
	svcsrv "auth/service"
	"auth/user"
	"context"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/sync/errgroup"
	stdgrpc "google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// create loggers
	l := kitlog.With(InitLogger(os.Stdout), "level", "info")
	errorLogger := kitlog.With(InitLogger(os.Stdout), "level", "error")
	l.Log("message", "Hello, I am alive")
	defer l.Log("message", "goodbye, good luck")

	// create DB
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		"Europe/Kiev",
	)))
	if err != nil {
		errorLogger.Log("db error:", err)
		os.Exit(1)
	}

	// DB migrate
	//tables, err := db.Migrator().GetTables()
	//for _, table := range tables {
	//	err = db.Debug().Migrator().DropTable(table)
	//	if err != nil {
	//		errorLogger.Log("rop table error:", err)
	//		os.Exit(1)
	//	}
	//}

	err = db.Debug().AutoMigrate(&types.User{}, &types.Account{}, &types.Service{}, &types.Permission{})
	if err != nil {
		errorLogger.Log("migrate error:", err)
		os.Exit(1)
	}

	// read keys
	//privateKey := os.Getenv("PRIVATE_KEY")
	//var pk *ecdsa.PrivateKey
	//if privateKey == "" {
	//	pk, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//} else {
	//	pk, err = ec.PEM2PrivateKey([]byte(privateKey))
	//}
	//if err != nil {
	//	errorLogger.Log("private key error:", err)
	//	os.Exit(1)
	//}

	// create services
	//authorizer := authz.New("active", "read", "write", "admin")
	//userSvc := user.NewLoggingMiddleware(logger.New("[ user ]\t"), user.New(db))
	//authzSvc := authz.NewLoggingMiddleware(logger.New("[ authz ]\t"), authz.NewService(db))
	permissionSvc := permission.New(db)
	svcSvc := svcsrv.New(logger.New("[ service ]\t"), db)
	accountSvc := account.NewLoggingMiddleware(logger.New("[ service ]\t"), account.New(db))
	userSvc := user.NewLoggingMiddleware(logger.New("[ user ]\t"), user.New(db))
	//jwtSvc := jwt.New(
	//	logger.New("[ jwt ]\t"),
	//	jwt2.SigningMethodES256,
	//	[]string{jwt2.SigningMethodES256.Name},
	//	jwt.AccessClaimsFactory,
	//	pk,
	//)

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return InterruptHandler(ctx)
	})
	mgmtSvc := mgmt.New(logger.New("[ mgmt ]\t"), db, svcSvc, accountSvc, permissionSvc, userSvc)
	mgmtSvc = mgmtservice.LoggingMiddleware(l)(mgmtSvc)
	mgmtSvc = mgmtservice.ErrorLoggingMiddleware(l)(mgmtSvc)
	mgmtSvc = mgmtservice.RecoveringMiddleware(l)(mgmtSvc)

	mgmtEndpoints := mgmttransport.Endpoints(mgmtSvc)

	grpcAddr := ":9091"
	// Start authgrpc server.
	g.Go(func() error {
		return ServeGRPC(ctx, &mgmtEndpoints, grpcAddr, kitlog.With(l, "transport", "GRPC"))
	})

	httpAddr := ":8081"
	// Start authhttp server.
	g.Go(func() error {
		return ServeHTTP(ctx, &mgmtEndpoints, httpAddr, kitlog.With(l, "transport", "HTTP"))
	})

	if err := g.Wait(); err != nil {
		l.Log("error", err)
	}
}

// InitLogger initialize go-kit JSON logger with timestamp and caller.
func InitLogger(writer io.Writer) kitlog.Logger {
	l := kitlog.NewLogfmtLogger(writer)
	l = kitlog.With(l, "@timestamp", kitlog.DefaultTimestampUTC)
	l = kitlog.With(l, "caller", kitlog.DefaultCaller)
	return l
}

// InterruptHandler handles first SIGINT and SIGTERM and returns it as error.
func InterruptHandler(ctx context.Context) error {
	interruptHandler := make(chan os.Signal, 1)
	signal.Notify(interruptHandler, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-interruptHandler:
		return fmt.Errorf("signal received: %v", sig.String())
	case <-ctx.Done():
		return errors.New("signal listener: context canceled")
	}
}

// ServeGRPC starts new GRPC server on address and sends first error to channel.
func ServeGRPC(ctx context.Context, endpoints *mgmttransport.EndpointsSet, addr string, logger kitlog.Logger) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// Here you can add middlewares for authgrpc server.
	server := mgmtgrpc.NewGRPCServer(endpoints)
	grpcServer := stdgrpc.NewServer()
	proto.RegisterMgmtServer(grpcServer, server)
	logger.Log("listen on", addr)
	ch := make(chan error)
	go func() {
		ch <- grpcServer.Serve(listener)
	}()
	select {
	case err := <-ch:
		return fmt.Errorf("authgrpc server: serve: %v", err)
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return errors.New("authgrpc server: context canceled")
	}
}

// ServeHTTP starts new HTTP server on address and sends first error to channel.
func ServeHTTP(ctx context.Context, endpoints *mgmttransport.EndpointsSet, addr string, logger kitlog.Logger) error {
	handler := mgmthttp.NewHTTPHandler(endpoints)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	logger.Log("listen on", addr)
	ch := make(chan error)
	go func() {
		ch <- httpServer.ListenAndServe()
	}()
	select {
	case err := <-ch:
		if err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("authhttp server: serve: %v", err)
	case <-ctx.Done():
		return httpServer.Shutdown(context.Background())
	}
}
