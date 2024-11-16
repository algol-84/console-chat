package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	envNameRefreshKey           = "JWT_REFRESH_TOKEN_KEY"
	envNameAccessKey            = "JWT_ACCESS_TOKEN_KEY"
	envNameRefreshKeyExpiration = "JWT_REFRESH_TOKEN_EXPIRATION_MINUTES"
	envNameAccessKeyExpiration  = "JWT_ACCESS_TOKEN_EXPIRATION_MINUTES"
)

// TokenConfig интерфейс получения настроек токена
type TokenConfig interface {
	RefreshToken() string
	AccessToken() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type tokenConfig struct {
	refreshToken           string
	accessToken            string
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

// NewTokenConfig конструктор получения настроек токена
func NewTokenConfig() (TokenConfig, error) {
	refreshToken := os.Getenv(envNameRefreshKey)
	if len(refreshToken) == 0 {
		return nil, errors.New("jwt refresh token not found")
	}

	accessToken := os.Getenv(envNameAccessKey)
	if len(refreshToken) == 0 {
		return nil, errors.New("jwt access token not found")
	}

	refreshTokenExpiration, err := strconv.Atoi(os.Getenv(envNameRefreshKeyExpiration))
	if err != nil {
		return nil, errors.New("jwt refresh token expiration time not found")
	}

	accessTokenExpiration, err := strconv.Atoi(os.Getenv(envNameAccessKeyExpiration))
	if err != nil {
		return nil, errors.New("jwt refresh token expiration time not found")
	}

	return &tokenConfig{
			refreshToken:           refreshToken,
			accessToken:            accessToken,
			refreshTokenExpiration: time.Duration(refreshTokenExpiration) * time.Minute,
			accessTokenExpiration:  time.Duration(accessTokenExpiration) * time.Minute,
		},
		nil
}

func (cfg *tokenConfig) RefreshToken() string {
	return cfg.refreshToken
}

func (cfg *tokenConfig) AccessToken() string {
	return cfg.accessToken
}

func (cfg *tokenConfig) RefreshTokenExpiration() time.Duration {
	return cfg.refreshTokenExpiration
}

func (cfg *tokenConfig) AccessTokenExpiration() time.Duration {
	return cfg.accessTokenExpiration
}
