package config

import "time"

type Config struct {
	DBUser               string        `env:"MYSQL_USER" env-default:"user" env-required:"1"`
	DBPassword           string        `env:"MYSQL_PASSWORD" env-default:"pass" env-required:"1"`
	DBName               string        `env:"MYSQL_DATABASE" env-default:"example" env-required:"1"`
	DBHost               string        `env:"MYSQL_HOST" env-default:"localhost" env-required:"1"`
	DBPort               string        `env:"MYSQL_PORT" env-default:"3306" env-required:"1"`
	TgBotToken           string        `env:"TG_BOT_TOKEN" env-required:"1"`
	TgUser               int           `env:"TG_USER_ID" env-required:"1"`
	RequestPeriod        time.Duration `env:"APP_REQUEST_PERIOD_SECONDS" env-required:"1"`
	TgNotificatorEnabled bool          `env:"APP_TG_NOTIFICATOR" env-required:"1"`
}
