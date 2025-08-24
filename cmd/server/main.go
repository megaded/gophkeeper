package main

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/manager"
	"gophkeeper/internal/server"
	"gophkeeper/internal/storage/fileStorage/minio"
	"gophkeeper/internal/storage/postgre"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.SetupLogger("Info")
	cfg := config.GetConfig()
	storage := postgre.NewStorage(&cfg)
	identity := identity.CreateIdentityProvider(&cfg)
	userManager := manager.CreateUserManager(storage, identity)
	fileStorage, err := minio.NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	binaryManager := manager.NewBinaryManager(fileStorage, storage)
	creditCardManager := manager.NewCreditCardManager(cfg, storage)
	crypter := manager.NewCryptoManager(cfg)
	credManager := manager.NewCredentialsManager(&crypter, storage)
	textManager := manager.NewTextManager(storage)
	server := server.NewServer(cfg, &storage, userManager, identity, binaryManager, creditCardManager, credManager, textManager)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()
	server.Start(ctx)
}
