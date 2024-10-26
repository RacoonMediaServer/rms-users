package service

import (
	"github.com/RacoonMediaServer/rms-users/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type authClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (a authClaims) Valid() error {
	return nil
}

func (s Service) GenerateAccessToken(userId string) (string, error) {
	claims := authClaims{UserID: userId}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)
	return token.SignedString([]byte(config.Config().Security.Key))
}
