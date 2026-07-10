package config

import "time"

type (
	BookingConfig struct {
		ServerConfig   `yaml:"server"`
		PostgresConfig `yaml:"postgres"`
		KafkaConfig    `yaml:"kafka"`
		LoggerConfig   `yaml:"logger"`
		TickerConfig   `yaml:"ticker"`
		TokensConfig   `yaml:"jwt_token"`
	}

	ServerConfig struct {
		Port            int           `yaml:"port"`
		Host            string        `yaml:"host"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		IdleTimeout     time.Duration `yaml:"idle_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}

	PostgresConfig struct {
		DbUrl       string        `env:"POSTGRES_URL"`
		MaxConns    int32         `yaml:"max_conns"`
		MaxIdle     int           `yaml:"max_idle"`
		MaxLifetime time.Duration `yaml:"max_lifetime"`
	}

	KafkaConfig struct {
		ConnectionType string        `yaml:"conn_type"`
		Host           string        `yaml:"host"`
		Port           int           `yaml:"port"`
		Topics         []string      `yaml:"topics"`
		Partitions     int           `yaml:"partitions"`
		Replications   int           `yaml:"replications"`
		Timeout        time.Duration `yaml:"timeout"`
	}

	LoggerConfig struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	}

	TickerConfig struct {
		SlotInterval time.Duration `yaml:"slot_interval"`
		TickInterval time.Duration `yaml:"tick_interval"`
	}

	TokensConfig struct {
		AccessTokenConfig  `yaml:"access"`
		RefreshTokenConfig `yaml:"refresh"`
		JwtSecret          string `env:"JWT_SECRET"`
	}

	AccessTokenConfig struct {
		IssuedBy  string        `yaml:"iss"`
		ExpiredAt time.Duration `yaml:"exp"`
	}

	RefreshTokenConfig struct {
		IssuedBy  string        `yaml:"iss"`
		ExpiredAt time.Duration `yaml:"exp"`
	}
)
