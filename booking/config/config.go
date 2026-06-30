package config

import "time"

type (
	BookingConfig struct {
		ServerConfig   `yaml:"server"`
		PostgresConfig `yaml:"postgres"`
		LoggerConfig   `yaml:"logger"`
		CronConfig     `yaml:"cron"`
	}

	ServerConfig struct {
		Port            int           `yaml:"port"`
		Host            string        `yaml:"host"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		IddleTimeout    time.Duration `yaml:"iddle_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}

	PostgresConfig struct {
		DB_URL      string        `env:"POSTGRES_URL"`
		MaxConns    int32         `yaml:"max_conns"`
		MaxIddle    int           `yaml:"max_iddle"`
		MaxLifetime time.Duration `yaml:"max_lifetime"`
	}

	LoggerConfig struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	}

	CronConfig struct {
		SlotInterval    time.Duration `yaml:"slot_interval"`
		JobStartHours   int           `yaml:"job_start_hours"`
		JobStartMinutes int           `yaml:"job_start_minutes"`
	}
)
