package config

import (
	"os"
)

type Config struct {
	APP         string
	Environment string // development, staging, production

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	Context struct {
		TimeOut string
	}

	Kafka struct {
		Brokers string
		GroupID string
		Topic   string
	}

	AWSS3 struct {
		AWSAccessKeyID     string
		AWSSecretAccessKey string
		BucketName         string
		Region             string
	}
	GinMode string // debug, test, release

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string

	RedisHost string
	RedisPort string

	SigningKey string
	AccessTTL  string
	RefreshTTL string

	CSVFilePath  string
	ConfFilePath string

	LogLevel string

	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
}

func Load() *Config {
	var cfg Config

	// general configuration
	cfg.APP = getEnv("APP", "mini-twitter")
	cfg.Environment = getEnv("ENVIRONMENT", "develop")
	cfg.LogLevel = getEnv("LOG_LEVEL", "debug")
	cfg.Context.TimeOut = getEnv("CONTEXT_TIMEOUT", "7s")
	cfg.GinMode = getEnv("GIN_MODE", "debug")

	// server configuration
	cfg.Server.Host = getEnv("SERVER_HOST", "your_host")
	cfg.Server.Port = getEnv("SERVER_PORT", ":your_port")
	cfg.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	cfg.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	cfg.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	cfg.PostgresHost = getEnv("POSTGRES_HOST", "postgres_host")
	cfg.PostgresPort = getEnv("POSTGRES_PORT", "postgres_port")
	cfg.PostgresUser = getEnv("POSTGRES_USER", "postgres_user")
	cfg.PostgresPassword = getEnv("POSTGRES_PASSWORD", "postgres-password")
	cfg.PostgresDatabase = getEnv("POSTGRES_DATABASE", "postgres_db")
	cfg.PostgresSSLMode = getEnv("POSTGRES_SSL_MODE", "disable")

	// AWS S3 config
	cfg.AWSS3.AWSAccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "your_access_key_id")
	cfg.AWSS3.AWSSecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "your_secret_access_key")
	cfg.AWSS3.BucketName = getEnv("AWS_BUCKET_NAME", "your_aws_s3_bucket")
	cfg.AWSS3.Region = getEnv("AWS_REGION", "your_region")

	// kafka configuration
	cfg.Kafka.Brokers = getEnv("KAFKA_BROKER", "kafka_broker")
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC", "kafka_topic_name")

	// redis configuration
	cfg.RedisHost = getEnv("REDIS_HOST", "redis_host")
	cfg.RedisPort = getEnv("REDIS_PORT", "redis_port")

	// token configuration
	cfg.SigningKey = getEnv("SIGNING_KEY", "secret_key")
	cfg.AccessTTL = getEnv("ACCESS_TTL", "6h")
	cfg.RefreshTTL = getEnv("REFRESH_TTL", "24h")

	cfg.CSVFilePath = getEnv("CSV_FILE_PATH", "path_to_csv")
	cfg.ConfFilePath = getEnv("CONF_FILE_PATH", "path_to_conf")

	cfg.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTPPort = getEnv("SMTP_PORT", "587")
	cfg.SMTPEmail = getEnv("SMTP_EMAIL", "your_email_address")
	cfg.SMTPPassword = getEnv("SMTP_PASSWORD", "your_email_password")

	return &cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
