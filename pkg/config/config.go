package config

import "os"

type Config struct {
	Host string
	Port string

	Token struct {
		Secret string
	}

	PostgresDB struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
}

func NewConfig() *Config {
	config := Config{}

	config.Host = getEnv("SERVER_HOST", "localhost")
	config.Port = getEnv("SERVER_PORT", ":8080")

	config.Token.Secret = getEnv("TOKEN_SECRET", "secret")

	config.PostgresDB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.PostgresDB.Port = getEnv("POSTGRES_PORT", "5432")
	config.PostgresDB.User = getEnv("POSTGRES_USER", "postgres")
	config.PostgresDB.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	config.PostgresDB.DBName = getEnv("POSTGRES_DB_NAME", "chat_app")
	config.PostgresDB.SSLMode = getEnv("POSTGRES_SSL_MODE", "disable")

	return &config
}

func getEnv(key, default_value string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return default_value
	}

	return value
}
