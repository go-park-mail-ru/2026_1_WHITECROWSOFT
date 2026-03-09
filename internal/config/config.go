package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	PORT             = "8000"
	SHUTDOWN_TIMEOUT = 5
)

type JWTConfig struct {
	Secret string
	Secure bool
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

type Config struct {
	JWT    JWTConfig
	Server ServerConfig
}

func Load() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET was not found, shutting down...")
	}

	secure := os.Getenv("IS_SECURE") == "true"

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	shutdownTimeout := SHUTDOWN_TIMEOUT * time.Second
	if timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			shutdownTimeout = time.Duration(timeout) * time.Second
		}
	}

	return &Config{
		JWT: JWTConfig{
			Secret: jwtSecret,
			Secure: secure,
		},
		Server: ServerConfig{
			Port:            port,
			ShutdownTimeout: shutdownTimeout,
		},
	}
}
