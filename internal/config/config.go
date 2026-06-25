package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config собирает все настройки приложения
type Config struct {
	DataBase DataBaseConfig
	Server   ServerConfig
}

// DataBaseConfig - настройки подключения к PostgreSQL
type DataBaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN собирает строку подключения к БД
func (c DataBaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// ServerConfig - настройки HTTP сервера
type ServerConfig struct {
	Port string
}

// getEnv читает строковую переменную окружения
// Если переменная не задана, то возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// getEnvAsInt читает числовую переменную окружения
func getEnvAsInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func Load() *Config {
	return &Config{
		DataBase: DataBaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "horror"),
			Password: getEnv("DB_PASSWORD", "horror_secret"),
			DBName:   getEnv("DB_NAME", "horror_facts"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}
