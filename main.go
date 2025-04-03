package main

import (
	"embed"
	"net/http"

	"github.com/wymli/xserver/view/handlers"
	"github.com/wymli/xserver/view/server"
	"github.com/wymli/xserver/view/server/view_server/gin"
)

//go:embed cmd/*
var staticFS embed.FS

func main() {
	app := &handlers.App{}

	ginServer := &gin.GinServer{Addr: "0.0.0.0:9999"}

	s := &server.Server{
		App:                 app,
		StaticFS:            http.FS(staticFS),
		StaticFSServingPath: "/fe",
	}

	s.RegisterViewServer(ginServer)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
