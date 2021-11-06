package config

import (
	"os"
	"path/filepath"
	"time"

	"async_arch/internal/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AppPort         string `yaml:"app_port"`
	AppHost         string `yaml:"app_host"`
	GoogleAppSecret string `yaml:"google_app_secret"`
	GoogleAppID     string `yaml:"google_app_id"`
	JWTKey          string `yaml:"jwt_key"`
	ConfigDB        DBConf `yaml:"conf_db"`
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

func InitConf(confFile string) (cfg *AppConfig) {
	logger.Info("Try read config from file", zap.String("path", confFile))
	file, err := os.Open(filepath.Clean(confFile))
	if err != nil {
		logger.Fatal("Error open config file", err)
	}
	defer func() {
		if e := file.Close(); e != nil {
			logger.Fatal("Error close config file", e)
		}
	}()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatal("Invalid config file", err)
	}

	logger.Info("Config ok")
	return
}
