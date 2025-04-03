package server

import (
	"fmt"
	"net/http"

	"github.com/wymli/xserver/view/handlers"
)

type ViewServerI interface {
	RegisterAppRoute(app *handlers.App)
	RegisterStaticFS(path string, fs http.FileSystem)
	Run() error
}

type Server struct {
	App                 *handlers.App
	StaticFS            http.FileSystem
	StaticFSServingPath string
	viewServers         []ViewServerI // maybe native_http_server, gin, hertz, kitex, grpc, ...
}

func (s *Server) RegisterViewServer(vs ViewServerI) {
	vs.RegisterAppRoute(s.App)
	vs.RegisterStaticFS(s.StaticFSServingPath, s.StaticFS)

	if s.viewServers == nil {
		s.viewServers = []ViewServerI{}
	}
	s.viewServers = append(s.viewServers, vs)
}

func (s *Server) Run() error {
	if len(s.viewServers) == 0 {
		return fmt.Errorf("no registered view server found")
	}

	errCh := make(chan error, len(s.viewServers))
	for _, v := range s.viewServers {
		go func() {
			errCh <- v.Run()
		}()
	}

	for e := range errCh {
		// server run fail if any view server fail
		if e != nil {
			return e
		}
	}

	return nil
}
