package api_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/RevittConsulting/mdbx-viewer/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Deps interface {
	GetRouter() chi.Router
}

type ApiServer struct {
	ctx    context.Context
	config *config.Config
	deps   Deps
}

func NewApiServer(ctx context.Context, config *config.Config, deps Deps) *ApiServer {
	return &ApiServer{
		ctx:    ctx,
		config: config,
		deps:   deps,
	}
}

func (s *ApiServer) Run() error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	// HTTP SERVER
	port := fmt.Sprintf(":%v", s.config.HttpPort)
	server := http.Server{
		Addr:    port,
		Handler: s.deps.GetRouter(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("error starting http: %v", err)
		}
	}()

	log.Println("server started on port " + port)

	// SHUTDOWN
	<-ctx.Done()
	log.Println("shutting down server")

	err := server.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("error shutting down the server: %v", err)
	}

	return nil
}
