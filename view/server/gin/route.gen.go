package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/wymli/xserver/view/handlers"
)

func RegisterAppRouteGen(engine *gin.Engine, app *handlers.App) {
	engine.Handle("GET", "/api/v1/ping", WrapGinHandler(handlers.WrapMiddlewares(app.Ping)))
	engine.Handle("GET", "/api/v1/echo", WrapGinHandler(handlers.WrapMiddlewares(app.Echo)))
	engine.Handle("POST", "/api/v1/user/login", WrapGinHandler(handlers.WrapMiddlewares(app.LoginUser)))
	engine.Handle("GET", "/api/v1/user", WrapGinHandler(handlers.WrapMiddlewares(app.GetUser)))
}
