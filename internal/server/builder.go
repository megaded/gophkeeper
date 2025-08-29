package server

import (
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/logger"
)

type ServerBuilder interface {
	SetConfig(cfg config.Config) ServerBuilder
	SetStorage(storage Storager) ServerBuilder
	SetUserManager(userManager UserManager) ServerBuilder
	SetIdentityProvider(identityProvider identity.IdentityProvider) ServerBuilder
	SetBinaryManager(binaryManager BinaryManager) ServerBuilder
	SetCreditCardManager(creditCardManager CreditCardManager) ServerBuilder
	SetTextManager(textManager TextManager) ServerBuilder
	Build() Server
}

func GetServerBuilder() ServerBuilder {
	return NewBuilderServer()
}

func (b BuilderServer) SetConfig(cfg config.Config) ServerBuilder {
	b.server.Cfg = cfg
	return b
}
func (b BuilderServer) SetTextManager(textManager TextManager) ServerBuilder {
	b.server.TextManager = textManager
	return b
}

func (b BuilderServer) SetStorage(storage Storager) ServerBuilder {
	b.server.Storage = storage
	return b
}

func (b BuilderServer) SetUserManager(userManager UserManager) ServerBuilder {
	b.server.UserManager = userManager
	return b
}
func (b BuilderServer) SetIdentityProvider(identityProvider identity.IdentityProvider) ServerBuilder {
	b.server.IdentityProvider = identityProvider
	return b
}

func (b BuilderServer) SetBinaryManager(binaryManager BinaryManager) ServerBuilder {
	b.server.BinaryManager = binaryManager
	return b
}

func (b BuilderServer) SetCreditCardManager(creditCardManager CreditCardManager) ServerBuilder {
	b.server.CreditCardManager = creditCardManager
	return b
}

func (b BuilderServer) SetCredentialsManager(creditCardManager CredentialsManager) ServerBuilder {
	b.server.CredManager = creditCardManager
	return b
}

func (b BuilderServer) Build() Server {
	return b.server
}

type BuilderServer struct {
	server Server
}

func NewBuilderServer() BuilderServer {
	logger.SetupLogger("Info")
	return BuilderServer{server: Server{}}
}
