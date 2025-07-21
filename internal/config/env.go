package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	AppEnv                  string `envconfig:"APP_ENV" default:"development"`
	AppPort                 string `envconfig:"APP_PORT" default:"8080"`
	AllowedOrigins          string `envconfig:"ALLOWED_ORIGINS" default:"*"`
	DatabaseURL             string `envconfig:"DATABASE_URL" required:"true"`
	RedisAddr               string `envconfig:"REDIS_ADDR" required:"true"`
	RedisPassword           string `envconfig:"REDIS_PASSWORD" required:"true"`
	MyAnimeListClientID     string `envconfig:"MYANIMELIST_CLIENT_ID" required:"true"`
	MyAnimeListClientSecret string `envconfig:"MYANIMELIST_CLIENT_SECRET" required:"true"`
	CloudinaryName          string `envconfig:"CLOUDINARY_NAME" required:"true"`
	CloudinaryAPIKey        string `envconfig:"CLOUDINARY_API_KEY" required:"true"`
	CloudinaryAPISecret     string `envconfig:"CLOUDINARY_API_SECRET" required:"true"`
}

func LoadEnv() (*Env, error) {
	_ = godotenv.Load(".env.local", ".env")

	var env Env
	if err := envconfig.Process("", &env); err != nil {
		return nil, err
	}

	return &env, nil
}
