package server

import (
	"context"
	"net/http"
	"time"

	"os"
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
	iface := os.Getenv("API_INTERFACE")
	if iface == "" {
		iface = "127.0.0.1"
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5144"
	}

	listenAddr := (iface + ":" + port)

	s.httpServer = &http.Server{
		Handler: s.API,
		Addr:    listenAddr,
	}

	go s.httpServer.ListenAndServe()
}

func (s *Server) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}
