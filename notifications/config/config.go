package config

import "time"

type NotificationsConfig struct {
	PostgresConfig `yaml:"postgres"`
	KafkaConfig    `yaml:"kafka"`
	LoggerConfig   `yaml:"logger"`
}

type PostgresConfig struct {
	DbUrl       string        `env:"POSTGRES_URL"`
	MaxConns    int32         `yaml:"max_conns"`
	MaxLifetime time.Duration `yaml:"max_lifetime"`
}

type KafkaConfig struct {
	ConnectionType string        `yaml:"conn_type"`
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	Topics         []string      `yaml:"topics"`
	Partitions     int           `yaml:"partitions"`
	Replications   int           `yaml:"replications"`
	Timeout        time.Duration `yaml:"timeout"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}
