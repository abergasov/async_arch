package main

import (
	"flag"

	"async_arch/internal/config"
	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/repository/exchanger"
	"async_arch/internal/repository/user"
	"async_arch/internal/routes/auth_routes"
	"async_arch/internal/service/auth"
	"async_arch/internal/storage/broker"
	"async_arch/internal/storage/database"

	"github.com/abergasov/schema_registry"

	"go.uber.org/zap"
)

var (
	appPrefix = "auth"
	confFile  = flag.String("config", "configs/auth.yml", "Config file path")
)

func main() {
	flag.Parse()
	logger.InitLogger(appPrefix)
	conf := config.InitConf(*confFile)
	conn := database.InitDBConnect(&conf.ConfigDB)

	userRepo := user.InitUserRepo(conn)
	brokerKfk := broker.InitKafkaProducer(&conf.ConfigBroker, entities.UserCUDBrokerTopic)
	brokerKfkBI := broker.InitKafkaProducer(&conf.ConfigBroker, entities.UserBIBrokerTopic)

	registry := schema_registry.InitRegistry([]int{1})
	userService := auth.InitUserService(userRepo, registry, brokerKfk, brokerKfkBI, conf.JWTKey)

	exc := exchanger.InitExchanger(conn) // exchange uuid to jwt
	router := auth_routes.InitAuthAppRouter(conf, userService, exc)
	logger.Info("start auth app", zap.String("url", conf.AppHost+":"+conf.AppPort))
	if err := router.InitRoutes(conf.JWTKey).Start(":" + conf.AppPort); err != nil {
		logger.Fatal("Common server error", err)
	}

}
