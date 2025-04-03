package handlers

import "github.com/wymli/xserver/view/model"

func (app *App) Ping(ctx *model.AppContext, req *model.PingReq) (rsp *model.PingRsp, err error) {
	return &model.PingRsp{}, nil
}
