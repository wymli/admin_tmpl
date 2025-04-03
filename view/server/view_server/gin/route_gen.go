package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/wymli/xserver/view/handlers"
)

func RegisterAppRouteGen(engine *gin.Engine, app *handlers.App) {
	engine.Handle("", "", WrapGinHandler(app.Echo))
}
