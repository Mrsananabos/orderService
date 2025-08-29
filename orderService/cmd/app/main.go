package main

import (
	"context"
	"log"
	_ "orderService/docs"
	"orderService/http/rest"
	"os"
	"os/signal"
	"syscall"
)

// @title Order Service
// @version 1.0
// @description Order service API Docs

// @host 	localhost:8080
// @BasePath /api
func main() {
	signals := []os.Signal{
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT,
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	server, err := rest.NewServer(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
