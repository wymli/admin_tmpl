package model

import "context"

type AppContext struct {
	ctx context.Context
}

type AppI interface {
	Ping(ctx *AppContext, req *PingReq) (rsp *PingRsp, err error) // file=ping method=get path=/api/v1/ping
	Echo(ctx *AppContext, req *EchoReq) (rsp *EchoRsp, err error) // file=ping method=get path=/api/v1/echo

	LoginUser(ctx *AppContext, req *LoginUserReq) (rsp *LoginUserRsp, err error) // file=user method=post path=/api/v1/user/login
	GetUser(ctx *AppContext, req *GetUserReq) (rsp *GetUserRsp, err error)       // file=user method=get path=/api/v1/user
}

type PingReq struct{}

type PingRsp struct{}

type EchoReq struct {
	Msg string `json:"msg,omitempty"`
}

type EchoRsp struct {
	Msg string `json:"msg,omitempty"`
}

type LoginUserReq struct {
	UserName string `form:"userName,omitempty"`
	Password string `form:"password,omitempty"`
}

type LoginUserRsp struct{}

type GetUserReq struct{}

type GetUserRsp struct {
	Name         string              `json:"name,omitempty"`
	Avatar       string              `json:"avatar,omitempty"`
	Email        string              `json:"email,omitempty"`
	Introduction string              `json:"introduction,omitempty"`
	Permissions  map[string][]string `json:"permissions,omitempty"`
}
