package handlers

import (
	"github.com/wymli/xserver/view/model"
)

func (app *App) Echo(ctx *model.AppContext, req *model.EchoReq) (rsp *model.EchoRsp, err error) {
	return &model.EchoRsp{}, nil
}
