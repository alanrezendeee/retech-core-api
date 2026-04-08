package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	envDatabaseURL = "DATABASE_URL"
	envHTTPAddr    = "HTTP_ADDR"
	envGinMode     = "GIN_MODE"
	envLogLevel    = "LOG_LEVEL"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	GinMode     string
	LogLevel    string
}

func Load() (*Config, error) {
	dbURL := strings.TrimSpace(os.Getenv(envDatabaseURL))
	if dbURL == "" {
		return nil, fmt.Errorf("%s is required", envDatabaseURL)
	}
	addr := strings.TrimSpace(os.Getenv(envHTTPAddr))
	if addr == "" {
		return nil, fmt.Errorf("%s is required", envHTTPAddr)
	}
	ginMode := strings.TrimSpace(os.Getenv(envGinMode))
	if ginMode == "" {
		return nil, fmt.Errorf("%s is required", envGinMode)
	}
	logLevel := strings.TrimSpace(os.Getenv(envLogLevel))
	if logLevel == "" {
		return nil, fmt.Errorf("%s is required", envLogLevel)
	}

	return &Config{
		DatabaseURL: dbURL,
		HTTPAddr:    addr,
		GinMode:     ginMode,
		LogLevel:    logLevel,
	}, nil
}
