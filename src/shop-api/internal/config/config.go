package config

import "os"

type Config struct {
	Port         string
	DSN          string
	AuthGRPCAddr string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DSN:          getEnv("DB_DSN", "postgres://be4:1@localhost:5432/shopapi?sslmode=disable"),
		AuthGRPCAddr: getEnv("AUTH_GRPC_ADDR", "localhost:50051"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
