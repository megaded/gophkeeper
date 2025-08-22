package identity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type IdentityProvider struct {
	key string
}

func CreateIdentityProvider(c *config.Config) IdentityProvider {
	return IdentityProvider{key: c.SecretKey}
}

func (id *IdentityProvider) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(id.key))
}

func (id *IdentityProvider) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(id.key), nil
	})

	if err != nil || !token.Valid {
		logger.Log.Error("failed to parse token", zap.Error(err))
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found")
	}

	return int(userIDFloat), nil
}

func (id *IdentityProvider) HashPassword(password string) string {
	hash := hmac.New(sha256.New, []byte(id.key))
	hash.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (id *IdentityProvider) VerifyPassword(hash string, password string) bool {
	hashedPassword := id.HashPassword(password)
	return hmac.Equal([]byte(hash), []byte(hashedPassword))
}
