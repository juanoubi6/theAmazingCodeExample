package config

import "os"

type Config struct {
	ENV        string
	PORT       string
	JWT_SECRET string
	CORS       string

	NATS_URL string

	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	RABBITMQ_HOST     string
	RABBITMQ_PORT     string

	DB_TYPE     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string

	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_BUCKET            string
	AWS_REGION            string

	SENDGRID_KEY_ID string

	GOOGLE_PLACES_API_KEY string
	GOOGLE_CLIENT_ID      string
	GOOGLE_CLIENT_SECRET  string

	TWILIO_SID        string
	TWILIO_AUTH_TOKEN string
	TWILIO_ACC_PHONE  string
}

var instance *Config

func GetConfig() *Config {
	if instance == nil {
		config := newConfig()
		instance = &config
	}
	return instance
}

func newConfig() Config {
	return Config{
		ENV:        GetEnv("ENV", "develop"),
		PORT:       GetEnv("PORT", "5000"),
		JWT_SECRET: GetEnv("JWT_SECRET", "j8Ah4kO3"),
		CORS:       GetEnv("CORS", ""),

		NATS_URL: GetEnv("NATS_URL", "0.0.0.0:4222"),

		RABBITMQ_HOST:     GetEnv("RABBITMQ_HOST", "localhost"),
		RABBITMQ_PORT:     GetEnv("RABBITMQ_PORT", "5672"),
		RABBITMQ_USER:     GetEnv("RABBITMQ_USER", "guest"),
		RABBITMQ_PASSWORD: GetEnv("RABBITMQ_PASSWORD", "guest"),

		DB_TYPE:     GetEnv("DB_TYPE", "mysql"),
		DB_USERNAME: GetEnv("DB_USERNAME", "root"),
		DB_PASSWORD: GetEnv("DB_PASSWORD", "root"),
		DB_HOST:     GetEnv("DB_HOST", "127.0.0.1"),
		DB_PORT:     GetEnv("DB_PORT", "3306"),
		DB_NAME:     GetEnv("DB_NAME", "amazing-code-database"),

		AWS_ACCESS_KEY_ID:     GetEnv("AWS_ACCESS_KEY_ID", ""),
		AWS_SECRET_ACCESS_KEY: GetEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWS_BUCKET:            GetEnv("AWS_BUCKET", ""),
		AWS_REGION:            GetEnv("AWS_REGION", ""),

		SENDGRID_KEY_ID: GetEnv("SENDGRID_KEY_ID", ""),

		GOOGLE_PLACES_API_KEY: GetEnv("GOOGLE_PLACES_API_KEY", ""),
		GOOGLE_CLIENT_ID:      GetEnv("GOOGLE_CLIENT_ID", ""),
		GOOGLE_CLIENT_SECRET:  GetEnv("GOOGLE_CLIENT_SECRET", ""),

		TWILIO_SID:        GetEnv("TWILIO_SID", ""),
		TWILIO_AUTH_TOKEN: GetEnv("TWILIO_AUTH_TOKEN", ""),
		TWILIO_ACC_PHONE:  GetEnv("TWILIO_ACC_PHONE", ""),
	}
}

func GetEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
