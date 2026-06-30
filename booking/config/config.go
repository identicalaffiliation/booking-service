package config

import "time"

type (
	BookingConfig struct {
		ServerConfig   `yaml:"server"`
		PostgresConfig `yaml:"postgres"`
		LoggerConfig   `yaml:"logger"`
		TickerConfig   `yaml:"ticker"`
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

	LoggerConfig struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	}

	TickerConfig struct {
		SlotInterval time.Duration `yaml:"slot_interval"`
		TickInterval time.Duration `yaml:"tick_interval"`
	}
)
