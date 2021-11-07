package main

import (
	"flag"

	"async_arch/internal/config"
	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/repository/user"
	"async_arch/internal/routes/task_routes"
	"async_arch/internal/service/task"
	"async_arch/internal/storage/broker"
	"async_arch/internal/storage/database"

	"go.uber.org/zap"
)

var (
	appPrefix = "task"
	confFile  = flag.String("config", "configs/task.yml", "Config file path")
)

func main() {
	flag.Parse()
	logger.InitLogger(appPrefix)
	conf := config.InitConf(*confFile)
	conn := database.InitDBConnect(&conf.ConfigDB)

	userRepo := user.InitUserRepo(conn)
	kfk := broker.InitKafkaConsumer(&conf.ConfigBroker, entities.UserCUDBrokerTopic)
	userService := task.InitUserTaskService(userRepo, kfk, conf.JWTKey)

	router := task_routes.InitAuthAppRouter(conf, userService)
	logger.Info("start auth app", zap.String("url", conf.AppHost+":"+conf.AppPort))
	if err := router.InitRoutes(conf.JWTKey).Start(":" + conf.AppPort); err != nil {
		logger.Fatal("Common server error", err)
	}
}
