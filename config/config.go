package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	ClickHouse  string
	NATSURL     string
	StreamName  string
	SubjectName string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", ":8080"),
		ClickHouse:  getEnv("CLICKHOUSE_URL", "http://localhost:8123"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		StreamName:  getEnv("NATS_STREAM", "items_stream"),
		SubjectName: getEnv("NATS_SUBJECT", "items"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
