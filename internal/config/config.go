package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	ConfigDB DBConf `yaml:"conf_db"`
}

type DBConf struct {
	Address        string        `yaml:"address"`
	Port           string        `yaml:"port"`
	User           string        `yaml:"user"`
	Pass           string        `yaml:"pass"`
	DBName         string        `yaml:"db_name"`
	MaxConnections int           `yaml:"max_connections"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
}

func InitConf(confFile string) (cfg *AppConfig, err error) {
	log.Println("Try read config from file", "path", confFile)
	file, err := os.Open(filepath.Clean(confFile))
	if err != nil {
		return nil, errors.New("Error open file:" + err.Error())
	}
	defer func() {
		if e := file.Close(); e != nil {
			log.Fatal("Error close config file", e)
		}
	}()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, errors.New("Invalid config file:" + err.Error())
	}

	log.Println("Config ok", "path", confFile)
	return
}
