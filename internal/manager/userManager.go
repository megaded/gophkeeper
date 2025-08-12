package manager

import (
	"context"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server"
)

type UserStorager interface {
	AddUser(ctx context.Context, login string, password string) error
}

type UserManager struct {
	storage  UserStorager
	identity identity.IdentityProvider
}

func CreateUserManager(s server.Storager) UserManager {
	return UserManager{storage: s}
}

func (u *UserManager) CreateUser(ctx context.Context, login string, password string) error {
	if login == "" || password == "" {
		return internal_error.ErrEmptyLoginOrPassword
	}
	hash := u.identity.HashPassword(password)
	return u.storage.AddUser(ctx, login, hash)
}
