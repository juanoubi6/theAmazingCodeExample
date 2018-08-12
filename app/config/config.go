package config

import "os"

type Config struct {
	ENV        string
	PORT       string
	JWT_SECRET string
	CORS       string

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

		SENDGRID_KEY_ID: GetEnv("SENDGRID_KEY_ID", "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"),

		GOOGLE_PLACES_API_KEY: GetEnv("GOOGLE_PLACES_API_KEY", "AIzaSyCv2CdoVDMQ6S8jz2vtDYGTAwJojnxHJus"),
		GOOGLE_CLIENT_ID:      GetEnv("GOOGLE_CLIENT_ID", "743009156834-jmfvt5p1uk2k1gmoqvakve9h4ru5aknj.apps.googleusercontent.com"),
		GOOGLE_CLIENT_SECRET:  GetEnv("GOOGLE_CLIENT_SECRET", "zJTGW-qyS7HHRvJt3TPtRH32"),

		TWILIO_SID:        GetEnv("TWILIO_SID", "ACb98ded914e1c12b0e276c4f164555f70"),
		TWILIO_AUTH_TOKEN: GetEnv("TWILIO_AUTH_TOKEN", "487d710fa087d1a2f5f757d579a0d367"),
		TWILIO_ACC_PHONE:  GetEnv("TWILIO_ACC_PHONE", "+19852418221"),
	}
}

func GetEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
