package main

import (
	"context"
	"github.com/RevittConsulting/mdbx-viewer/config"
	"github.com/RevittConsulting/mdbx-viewer/internal/api_server"
	"github.com/RevittConsulting/mdbx-viewer/internal/dependencies"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.InitializeConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v", err)
		return
	}

	deps := dependencies.NewDependencies(cfg)

	as := api_server.NewApiServer(ctx, cfg, deps)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()

	if err = as.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
