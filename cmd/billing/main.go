package main

import (
	"flag"

	"async_arch/internal/config"
	"async_arch/internal/entities"
	"async_arch/internal/logger"
	tRepo "async_arch/internal/repository/task"
	"async_arch/internal/repository/user"
	"async_arch/internal/routes/billing_routes"
	"async_arch/internal/service"
	"async_arch/internal/storage/broker"
	"async_arch/internal/storage/database"

	"go.uber.org/zap"

	"github.com/abergasov/schema_registry"
)

var (
	appPrefix = "task"
	confFile  = flag.String("config", "configs/billing.yml", "Config file path")
)

func main() {
	flag.Parse()
	logger.InitLogger(appPrefix)
	conf := config.InitConf(*confFile)
	conn := database.InitDBConnect(&conf.ConfigDB)

	userRepo := user.InitUserRepo(conn)
	taskRepo := tRepo.InitTaskRepo(conn)

	kfk := broker.InitKafkaConsumer(&conf.ConfigBroker, "billing", entities.UserCUDBrokerTopic)
	registry := schema_registry.InitRegistry([]int{1})
	service.InitUserReplicatorService(userRepo, registry, kfk)

	taskRegistry := schema_registry.InitRegistry([]int{2})
	kfkTask := broker.InitKafkaConsumer(&conf.ConfigBroker, "billing", entities.TaskCUDBrokerTopic)
	kfkTaskBI := broker.InitKafkaConsumer(&conf.ConfigBroker, "billing", entities.TaskBIBrokerTopic)
	service.InitTaskReplicatorService(taskRepo, taskRegistry, kfkTask, kfkTaskBI)

	router := billing_routes.InitBillingAppRouter(conf)
	logger.Info("start auth app", zap.String("url", conf.AppHost+":"+conf.AppPort))
	if err := router.InitRoutes(conf.JWTKey).Start(":" + conf.AppPort); err != nil {
		logger.Fatal("Common server error", err)
	}
}
