package main

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/manager"
	"gophkeeper/internal/server"
	"gophkeeper/internal/storage"
	"gophkeeper/internal/storage/fileStorage/minio"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.SetupLogger("Info")
	cfg := config.GetConfig()
	storage := storage.NewPgStorage(&cfg)
	userManager := manager.CreateUserManager(storage)
	fileStorage, err := minio.NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	server := server.NewServer(cfg, &storage, userManager, identity.CreateIdentityProvider(&cfg), fileStorage)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()
	server.Start(ctx)
}
