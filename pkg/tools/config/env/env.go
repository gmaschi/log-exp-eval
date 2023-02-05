package env

import (
	"os"
)

const (
	dbDriverKey         = "DB_DRIVER"
	dbSourceKey         = "DB_SOURCE"
	serverAddressKey    = "SERVER_ADDRESS"
	dbPortKey           = "POSTGRES_PORT"
	dbHostKey           = "DB_HOST"
	postgresUserKey     = "POSTGRES_USER"
	postgresPasswordKey = "POSTGRES_PASSWORD"
	postgresDBKey       = "POSTGRES_DB"
)

type Config struct {
	DbDriver         string `json:"DB_DRIVER"`
	DbSource         string `json:"DB_SOURCE"`
	ServerAddress    string `json:"SERVER_ADDRESS"`
	DBPort           string `json:"POSTGRES_PORT"`
	DBHost           string `json:"DB_HOST"`
	PostgresUser     string `json:"POSTGRES_USER"`
	PostgresPassword string `json:"POSTGRES_PASSWORD"`
	PostgresDB       string `json:"POSTGRES_DB"`
}

// NewConfig returns the config struct loaded with the environment variables.
func NewConfig() (Config, error) {
	return Config{
		DbDriver:         os.Getenv(dbDriverKey),
		DbSource:         os.Getenv(dbSourceKey),
		ServerAddress:    os.Getenv(serverAddressKey),
		DBPort:           os.Getenv(dbPortKey),
		DBHost:           os.Getenv(dbHostKey),
		PostgresUser:     os.Getenv(postgresUserKey),
		PostgresPassword: os.Getenv(postgresPasswordKey),
		PostgresDB:       os.Getenv(postgresDBKey),
	}, nil
}
