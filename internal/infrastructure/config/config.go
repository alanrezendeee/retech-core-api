package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Variáveis de ambiente obrigatórias. Load falha se qualquer uma estiver ausente
// ou for só espaços — o processo não deve subir sem todas definidas.
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

// Load lê e valida a configuração. Todas as constantes *env* acima precisam estar
// definidas no ambiente (ou vindas do .env carregado antes da chamada).
func Load() (*Config, error) {
	log.Println("🔍 Carregando configurações obrigatórias...")

	dbURL := strings.TrimSpace(os.Getenv(envDatabaseURL))
	addr := strings.TrimSpace(os.Getenv(envHTTPAddr))
	ginMode := strings.TrimSpace(os.Getenv(envGinMode))
	logLevel := strings.TrimSpace(os.Getenv(envLogLevel))

	var missing []string
	if dbURL == "" {
		missing = append(missing, envDatabaseURL)
	}
	if addr == "" {
		missing = append(missing, envHTTPAddr)
	}
	if ginMode == "" {
		missing = append(missing, envGinMode)
	}
	if logLevel == "" {
		missing = append(missing, envLogLevel)
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("required environment variables missing or empty: %s", strings.Join(missing, ", "))
	}

	log.Println("✅ Todas as configurações carregadas com sucesso!")

	return &Config{
		DatabaseURL: dbURL,
		HTTPAddr:    addr,
		GinMode:     ginMode,
		LogLevel:    logLevel,
	}, nil
}

// RequiredEnvKeys retorna os nomes das variáveis obrigatórias.
func RequiredEnvKeys() []string {
	return []string{envDatabaseURL, envHTTPAddr, envGinMode, envLogLevel}
}
