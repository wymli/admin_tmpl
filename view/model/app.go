package model

import "context"

type AppContext struct {
	ctx context.Context
}

type AppI interface {
	Ping(ctx *AppContext, req *PingReq) (rsp *PingRsp, err error) // method=get path=/api/v1/ping
	Echo(ctx *AppContext, req *EchoReq) (rsp *EchoRsp, err error) // method=get path=/api/v1/echo
}

