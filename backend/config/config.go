package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort   string
	HTTPPort   string
	ServerAddr string
	DBUrl      string
}

func LoadConfig() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found")
	}
	return &Config{
		GRPCPort:   getEvn("GRPC_PORT", "50051"),
		HTTPPort:   getEvn("HTTP_PORT", "8080"),
		ServerAddr: getEvn("SERVER_ADDRESS", "localhost:50051"),
		DBUrl:      buildDBUrl(),
	}
}

func getEvn(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func buildDBUrl() string {
	host := getEvn("DB_HOST", "localhost")
	user := getEvn("DB_USER", "postgres")
	password := getEvn("DB_PASSWORD", "postgres")
	name := getEvn("DB_NAME", "sysmonitor")
	port := getEvn("DB_PORT", "5432")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
}
