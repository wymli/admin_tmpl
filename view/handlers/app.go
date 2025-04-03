package handlers

import "github.com/wymli/xserver/view/model"

type App struct{}

type AppHandlerFn[Req, Rsp any] func(ctx *model.AppContext, req Req) (rsp Rsp, err error)

func PrintLog[Req, Rsp any](f AppHandlerFn[Req, Rsp]) AppHandlerFn[Req, Rsp] {
	return func(ctx *model.AppContext, req Req) (rsp Rsp, err error) {
		return f(ctx, req)
	}
}
