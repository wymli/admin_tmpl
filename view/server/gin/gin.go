package gin

import (
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wymli/xserver/utils/gresult"
	"github.com/wymli/xserver/view/handlers"
	"github.com/wymli/xserver/view/model"
)

type GinServer struct {
	Addr string

	engine *gin.Engine
}

func NewGinServer(addr string) *GinServer {
	return &GinServer{
		Addr:   addr,
		engine: gin.Default(),
	}
}

type ResponseObj struct {
	Code    int    `json:"code,omitempty"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func WrapGinHandler[Req, Rsp model.ViewModelI](h handlers.AppHandlerFn[Req, Rsp]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ResponseObj{Code: -1, Data: nil, Message: "bind failed"})
			return
		}

		rspData, err := h(&model.AppContext{}, &req)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ResponseObj{Code: -1, Data: nil, Message: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, ResponseObj{Code: 0, Data: rspData, Message: "ok"})
	}
}

func (s *GinServer) RegisterAppRoute(app *handlers.App) {
	RegisterAppRouteGen(s.engine, app)
}

func (s *GinServer) RegisterStaticFS(staticFs fs.FS) {
	s.engine.StaticFS("/assets", http.FS(gresult.Must(fs.Sub(staticFs, "assets"))))
	s.engine.NoRoute(func(ctx *gin.Context) {
		// issue: https://stackoverflow.com/questions/43527073/golang-static-stop-index-html-redirection
		// ctx.FileFromFS("index.html", http.FS(gresult.Must(fs.Sub(staticFs, "fe/dist"))))

		// 未注册的api路由返回失败
		if strings.HasPrefix(ctx.Request.URL.Path, "/api") {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 其他路由返回前端主页
		f, _ := staticFs.Open("index.html")
		ctx.Status(200)
		_, _ = io.Copy(ctx.Writer, f)
	})
}

func (s *GinServer) Run() error {
	return s.engine.Run(s.Addr)
}
