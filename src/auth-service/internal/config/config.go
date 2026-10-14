package config

import "os"

type Config struct {
	GRPCPort  string
	DSN       string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		DSN:       getEnv("DB_DSN", "postgres://be4:1@localhost:5432/auth_db?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-prod"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
