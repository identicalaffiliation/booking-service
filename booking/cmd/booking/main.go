package main

import (
	"flag"

	"github.com/identicalaffiliation/booking-service/booking/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	_ = config.MustLoad(configPath)
}
