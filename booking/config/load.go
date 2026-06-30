package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad(configPath string) *BookingConfig {
	cfg := new(BookingConfig)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}

	return cfg
}
