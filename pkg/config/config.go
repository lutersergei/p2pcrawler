package config

import "time"

type Config struct {
	TgBotToken    string        `env:"TG_BOT_TOKEN" env-required:"1"`
	TgUser        int           `env:"TG_USER_ID" env-required:"1"`
	RequestPeriod time.Duration `env:"APP_REQUEST_PERIOD_SECONDS" env-required:"1"`
	SentryDSN     string        `env:"SENTRY_DSN" env-required:"1"`
}
