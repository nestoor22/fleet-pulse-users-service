package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type AuthConfig struct {
	JwtSecret                     string
	JwtAccessTokenExpireInMinutes int
	JwtRefreshTokenExpireInHours  int
	InviteSecret                  string
	InviteExpireInMinutes         int
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	jwtAccessTokenExpire, err := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_EXPIRE_IN_MINUTES", "10"))
	if err != nil {
		log.Fatal(err)
	}

	jwtRefreshTokenExpire, err := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_EXPIRE_IN_HOURS", "72"))
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", ":8000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		Auth: AuthConfig{
			JwtSecret:                     getEnv("JWT_SECRET", ""),
			JwtAccessTokenExpireInMinutes: jwtAccessTokenExpire,
			JwtRefreshTokenExpireInHours:  jwtRefreshTokenExpire,
			InviteSecret:                  getEnv("INVITE_SECRET", ""),
			InviteExpireInMinutes:         60,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
