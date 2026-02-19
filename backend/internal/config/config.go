package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	ServerPort       string
	JWTSecret        string
	JWTExpiryHours   int
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
	WebhookRateLimit int
}

func Load() *Config {
	if err := LoadEnv(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	rateLimit, _ := strconv.Atoi(getEnv("WEBHOOK_RATE_LIMIT", "100"))

	// Use 6379 as default Redis port
	redisPort := getEnv("REDIS_PORT", "6379")

	return &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "change-this-in-production"),
		JWTExpiryHours:   jwtExpiry,
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "root"),
		DBName:           getEnv("DB_NAME", "support_platform_db"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        redisPort,
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          redisDB,
		WebhookRateLimit: rateLimit,
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func LoadEnv() error {
	return godotenv.Load()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
