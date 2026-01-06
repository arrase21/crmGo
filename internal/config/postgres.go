package config

import (
	"log"
	"os"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func LoadPostgres() *PostgresConfig {
	return &PostgresConfig{
		Host:     mustGetEnv("DB_HOST"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     mustGetEnv("DB_USER"),
		Password: mustGetEnv("DB_PASSWORD"),
		DBName:   mustGetEnv("DB_NAME"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		TimeZone: getEnv("DB_TIMEZONE", "America/Bogota"),
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("")
	}
	return v
}
