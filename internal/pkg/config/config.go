package config

import (
	"github.com/spf13/cast"
	"os"
)

type Config struct {
	APP         string
	Environment string // development, staging, production
	HTTPHost    string
	HTTPPort    int
	CtxTimeout  string
	GinMode     string // debug, test, release

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string

	RedisHost string
	RedisPort int

	Timeout string

	SigningKey string
	AccessTTL  string
	RefreshTTL string

	CSVFilePath  string
	ConfFilePath string

	LogLevel string

	SMTPHost     string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string
}

func Load() *Config {
	var cfg Config

	cfg.APP = cast.ToString(getEnv("APP", "app"))
	cfg.Environment = cast.ToString(getEnv("ENVIRONMENT", "development"))
	cfg.HTTPHost = cast.ToString(getEnv("HTTP_HOST", "twitter"))
	cfg.HTTPPort = cast.ToInt(getEnv("HTTP_PORT", "7777"))
	cfg.CtxTimeout = cast.ToString(getEnv("CONTEXT_TIMEOUT", "7s"))
	cfg.GinMode = cast.ToString(getEnv("GIN_MODE", "debug"))

	cfg.PostgresHost = cast.ToString(getEnv("POSTGRES_HOST", "twitter_postgres"))
	cfg.PostgresPort = cast.ToString(getEnv("POSTGRES_PORT", "5432"))
	cfg.PostgresUser = cast.ToString(getEnv("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getEnv("POSTGRES_PASSWORD", "root"))
	cfg.PostgresDatabase = cast.ToString(getEnv("POSTGRES_DATABASE", "twitter_db"))
	cfg.PostgresSSLMode = cast.ToString(getEnv("POSTGRES_SSL_MODE", "disable"))

	cfg.RedisHost = cast.ToString(getEnv("REDIS_HOST", "twitter_redis"))
	cfg.RedisPort = cast.ToInt(getEnv("REDIS_PORT", "6379"))

	cfg.Timeout = cast.ToString(getEnv("CONTEXT_TIMEOUT", "7s"))

	cfg.SigningKey = cast.ToString(getEnv("SIGNING_KEY", "twitter-secret"))
	cfg.AccessTTL = cast.ToString(getEnv("ACCESS_TTL", "6h"))
	cfg.RefreshTTL = cast.ToString(getEnv("REFRESH_TTL", "24h"))

	cfg.CSVFilePath = cast.ToString(getEnv("CSV_FILE_PATH", "./config/auth.csv"))
	cfg.ConfFilePath = cast.ToString(getEnv("CONF_FILE_PATH", "./config/auth.conf"))

	cfg.SMTPHost = cast.ToString(getEnv("SMTP_HOST", "smtp.gmail.com"))
	cfg.SMTPPort = cast.ToInt(getEnv("SMTP_PORT", "587"))
	cfg.SMTPEmail = cast.ToString(getEnv("SMTP_EMAIL", "xasannosirov094@gmail.com"))
	cfg.SMTPPassword = cast.ToString(getEnv("SMTP_PASSWORD", "zvhkpndjwrkiemci"))

	return &cfg
}

func getEnv(key string, defaultVal interface{}) interface{} {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
