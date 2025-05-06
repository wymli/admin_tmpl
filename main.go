package main

import (
	"embed"
	"io/fs"

	"github.com/wymli/xserver/utils/gresult"
	"github.com/wymli/xserver/view/handlers"
	"github.com/wymli/xserver/view/server"
	"github.com/wymli/xserver/view/server/gin"
)

//go:embed fe/dist
var staticFS embed.FS

func main() {
	app := &handlers.App{}
	s := &server.Server{}

	{
		ginServer := gin.NewGinServer("0.0.0.0:9999")
		ginServer.RegisterAppRoute(app)
		ginServer.RegisterStaticFS(gresult.Must(fs.Sub(staticFS, "fe/dist")))
		s.RegisterViewServer(ginServer)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}
