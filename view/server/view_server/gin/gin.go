package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wymli/xserver/view/handlers"
	"github.com/wymli/xserver/view/model"
)

type GinServer struct {
	Addr string

	engine *gin.Engine
}

type ResponseObj struct {
	Code    int
	Data    any
	Message string
}

func WrapGinHandler[Req, Rsp any](h handlers.AppHandlerFn[Req, Rsp]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ResponseObj{Code: -1, Data: nil, Message: "bind failed"})
			return
		}

		rspData, err := h(&model.AppContext{}, req)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseObj{Code: -1, Data: nil, Message: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, ResponseObj{Code: 0, Data: rspData, Message: "ok"})
		return
	}
}

func (s *GinServer) RegisterAppRoute(app *handlers.App) {
	RegisterAppRouteGen(s.engine, app)
}

func (s *GinServer) RegisterStaticFS(path string, fs http.FileSystem) {
	s.engine.StaticFS(path, fs)
}

func (s *GinServer) Run() error {
	return s.engine.Run(s.Addr)
}
