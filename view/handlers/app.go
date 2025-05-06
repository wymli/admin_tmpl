package handlers

import "github.com/wymli/xserver/view/model"

type App struct{}

type AppHandlerFn[Req, Rsp model.ViewModelI] func(ctx *model.AppContext, req *Req) (rsp *Rsp, err error)

func WrapMiddlewares[Req, Rsp model.ViewModelI](f AppHandlerFn[Req, Rsp]) AppHandlerFn[Req, Rsp] {
	return PrintLog(f)
}

func PrintLog[Req, Rsp model.ViewModelI](f AppHandlerFn[Req, Rsp]) AppHandlerFn[Req, Rsp] {
	return func(ctx *model.AppContext, req *Req) (rsp *Rsp, err error) {
		return f(ctx, req)
	}
}
