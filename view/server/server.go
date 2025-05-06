package server

import (
	"fmt"
)

type ViewServerI interface {
	Run() error
}

type Server struct {
	viewServers []ViewServerI // maybe native_http_server, gin, hertz, kitex, grpc, ...
}

func (s *Server) RegisterViewServer(vs ViewServerI) {
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
