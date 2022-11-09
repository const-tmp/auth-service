// Code generated by microgen 1.0.5. DO NOT EDIT.

package transport

import (
	auth "auth/auth"
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
)

func Endpoints(svc auth.Service) EndpointsSet {
	return EndpointsSet{
		LoginEndpoint:    LoginEndpoint(svc),
		RegisterEndpoint: RegisterEndpoint(svc),
	}
}

func RegisterEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(arg0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*RegisterRequest)
		res0, res1 := svc.Register(arg0, req.Login, req.Password, req.Service, req.AccountID)
		return &RegisterResponse{Ok: res0}, res1
	}
}

func LoginEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(arg0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*LoginRequest)
		res0, res1 := svc.Login(arg0, req.Login, req.Password, req.Service)
		return &LoginResponse{At: res0}, res1
	}
}