package main

import (
	"flag"

	"async_arch/internal/config"
)

var confFile = flag.String("config", "configs/auth.yml", "Config file path")

func main() {
	flag.Parse()
	conf, err := config.InitConf(*confFile)
	println(err, conf)
}
