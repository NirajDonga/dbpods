package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port      string
	Namespace string
}

func LoadConfig() *AppConfig {
	_ = godotenv.Load()

	return &AppConfig{
		Port:      mustGetEnv("PORT"),
		Namespace: mustGetEnv("K8S_NAMESPACE"),
	}
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists || value == "" {
		log.Fatalf("CRITICAL STARTUP ERROR: Environment variable '%s' is required but not set.", key)
	}

	return value
}
