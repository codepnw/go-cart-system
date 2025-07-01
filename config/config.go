package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppPort       string
	DBAddr        string
	JWTSecretKey  string
	JWTRefreshKey string
}

func InitEnvConfig(envPath string) (*EnvConfig, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	env := &EnvConfig{}

	env.AppPort = envGetString("APP_PORT", ":8080")
	env.DBAddr = envGetString("DB_ADDR", "")
	env.JWTSecretKey = envGetString("JWT_SECRET_KEY", "")
	env.JWTRefreshKey = envGetString("JWT_REFRESH_KEY", "")

	if err := validate(env); err != nil {
		return nil, err
	}

	return env, nil
}

func validate(config *EnvConfig) error {
	if config.DBAddr == "" {
		return errors.New("DB_ADDR is required")
	}
	if config.JWTSecretKey == "" {
		return errors.New("JWT_SECRET_KEY is required")
	}
	if config.JWTRefreshKey == "" {
		return errors.New("JWT_REFRESH_KEY is required")
	}
	return nil
}

func envGetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
