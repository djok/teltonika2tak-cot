package config

import (
	"context"

	"github.com/sirupsen/logrus"
)

type TeltonikaConfig struct {
	Host         string
	Port         int
	AllowedIMEIs []string
}

type Config struct {
	log             *logrus.Logger
	teltonikaConfig *TeltonikaConfig
}

func NewConfig(log *logrus.Logger, teltonikaConfig *TeltonikaConfig) *Config {
	return &Config{
		log:             log,
		teltonikaConfig: teltonikaConfig,
	}
}

func (c *Config) GetTeltonikaConfig() *TeltonikaConfig {
	return c.teltonikaConfig
}

func (c *Config) GetLogger() *logrus.Logger {
	return c.log
}

func GetLogger(ctx context.Context) *logrus.Logger {
	config := ctx.Value(ContextConfigKey).(*Config)
	return config.GetLogger()
}
