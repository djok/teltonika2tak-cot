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

type TakConfig struct {
	Host     string
	Port     int
	Protocol string
}

type Config struct {
	log             *logrus.Logger
	teltonikaConfig *TeltonikaConfig
	takConfig       *TakConfig
}

func NewConfig(log *logrus.Logger, teltonikaConfig *TeltonikaConfig, takConfig *TakConfig) *Config {
	return &Config{
		log:             log,
		teltonikaConfig: teltonikaConfig,
		takConfig:       takConfig,
	}
}

func (c *Config) GetTeltonikaConfig() *TeltonikaConfig {
	return c.teltonikaConfig
}

func (c *Config) GetTakConfig() *TakConfig {
	return c.takConfig
}

func (c *Config) GetLogger() *logrus.Logger {
	return c.log
}

func GetLogger(ctx context.Context) *logrus.Logger {
	config := ctx.Value(ContextConfigKey).(*Config)
	return config.GetLogger()
}
