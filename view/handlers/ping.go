package handlers

import "github.com/wymli/xserver/view/model"

func (app *App) Ping(ctx *model.AppContext, req *model.PingReq) (rsp *model.PingRsp, err error) {
	return &model.PingRsp{}, nil
}

func (app *App) Echo(ctx *model.AppContext, req *model.EchoReq) (rsp *model.EchoRsp, err error) {
	return &model.EchoRsp{Msg: req.Msg}, nil
}
