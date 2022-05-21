package server

import (
	"context"
	"net/http"
	"time"

	"reichard.io/libgen-opds/api"
)

type Server struct {
	API        *api.API
	httpServer *http.Server
}

func NewServer() *Server {
	api := api.NewApi()

	return &Server{
		API: api,
	}
}

func (s *Server) StartServer() {
	listenAddr := ("127.0.0.1:5144")

	s.httpServer = &http.Server{
		Handler: s.API.Router,
		Addr:    listenAddr,
	}

	go s.httpServer.ListenAndServe()
}

func (s *Server) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}
