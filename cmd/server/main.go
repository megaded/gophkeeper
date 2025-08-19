package main

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/manager"
	"gophkeeper/internal/server"
	"gophkeeper/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/* var filePath = "I:/Torrents/God Is a Bullet (2023)WEB-DLRip-AVC.mkv"
	os.ReadFile()
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	file.Read() */
	logger.SetupLogger("Info")
	cfg := config.GetConfig()
	storage := storage.NewPgStorage(&cfg)
	userManager := manager.CreateUserManager(storage)
	server := server.NewServer(cfg, &storage, userManager, identity.CreateIdentityProvider(&cfg))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()
	server.Start(ctx)
}
