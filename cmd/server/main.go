package main

import (
	"context"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server"
	"os"
	"os/signal"
	"syscall"
)

var (
	Version string
	Time    string
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		logger.Log.Warn(err.Error())
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()
	server.Start(ctx)
}
