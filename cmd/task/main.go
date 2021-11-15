package main

import (
	"flag"

	"async_arch/internal/config"
	"async_arch/internal/entities"
	"async_arch/internal/logger"
	tRepo "async_arch/internal/repository/task"
	"async_arch/internal/repository/user"
	"async_arch/internal/routes/task_routes"
	"async_arch/internal/service"
	"async_arch/internal/service/task"
	"async_arch/internal/storage/broker"
	"async_arch/internal/storage/database"

	"github.com/abergasov/schema_registry"

	"go.uber.org/zap"
)

const appPrefix = "task"

var (
	confFile = flag.String("config", "configs/task.yml", "Config file path")
)

func main() {
	flag.Parse()
	logger.InitLogger(appPrefix)
	conf := config.InitConf(*confFile)
	conn := database.InitDBConnect(&conf.ConfigDB)

	userRepo := user.InitUserRepo(conn)
	taskRepo := tRepo.InitTaskRepo(conn)

	service.InitUserReplicatorService(
		userRepo,
		schema_registry.InitRegistry([]int{1}),
		broker.InitKafkaConsumer(&conf.ConfigBroker, appPrefix, entities.UserCUDBrokerTopic),
	)

	brokerKfk := broker.InitKafkaProducer(&conf.ConfigBroker, entities.TaskCUDBrokerTopic)
	registryTask := schema_registry.InitRegistry([]int{2})
	brokerKfkBI := broker.InitKafkaProducer(&conf.ConfigBroker, entities.TaskBIBrokerTopic)
	taskService := task.InitTaskManager(registryTask, taskRepo, userRepo, brokerKfk, brokerKfkBI)

	router := task_routes.InitAuthAppRouter(conf, taskService)
	logger.Info("start auth app", zap.String("url", conf.AppHost+":"+conf.AppPort))
	if err := router.InitRoutes(conf.JWTKey).Start(":" + conf.AppPort); err != nil {
		logger.Fatal("Common server error", err)
	}
}
