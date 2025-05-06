package handlers

import (
	"fmt"

	"github.com/wymli/xserver/utils"
	"github.com/wymli/xserver/view/model"
)

func (app *App) LoginUser(ctx *model.AppContext, req *model.LoginUserReq) (rsp *model.LoginUserRsp, err error) {
	fmt.Printf("LoginUser: %s\n", utils.Json(req))
	return &model.LoginUserRsp{}, nil
}

func (app *App) GetUser(ctx *model.AppContext, req *model.GetUserReq) (rsp *model.GetUserRsp, err error) {
	name := "hello_world"
	return &model.GetUserRsp{
		Name:   name,
		Avatar: fmt.Sprintf("//api.randomx.ai/avatar/%s", name),
	}, nil
}
