package main

import (
	"flag"

	"async_arch/internal/config"
	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/repository/account"
	tRepo "async_arch/internal/repository/task"
	"async_arch/internal/repository/user"
	"async_arch/internal/routes/billing_routes"
	"async_arch/internal/service"
	"async_arch/internal/service/billing"
	"async_arch/internal/storage/broker"
	"async_arch/internal/storage/database"

	"go.uber.org/zap"

	"github.com/abergasov/schema_registry"
)

const appPrefix = "billing"

var (
	confFile = flag.String("config", "configs/billing.yml", "Config file path")
)

func main() {
	flag.Parse()
	logger.InitLogger(appPrefix)
	conf := config.InitConf(*confFile)
	conn := database.InitDBConnect(&conf.ConfigDB)

	service.InitUserReplicatorService(
		user.InitUserRepo(conn),
		schema_registry.InitRegistry([]int{1}),
		broker.InitKafkaConsumer(&conf.ConfigBroker, appPrefix, entities.UserCUDBrokerTopic),
	)

	service.InitTaskReplicatorService(
		tRepo.InitTaskRepo(conn),
		schema_registry.InitRegistry([]int{2}),
		broker.InitKafkaConsumer(&conf.ConfigBroker, appPrefix, entities.TaskCUDBrokerTopic),
		broker.InitKafkaConsumer(&conf.ConfigBroker, appPrefix, entities.TaskBIBrokerTopic),
	)

	billing.InitAccounter(
		account.InitAccountRepo(conn),
		broker.InitKafkaConsumer(&conf.ConfigBroker, appPrefix, entities.UserBIBrokerTopic),
		broker.InitKafkaConsumer(&conf.ConfigBroker, appPrefix, entities.TaskBIBrokerTopic),
	)

	router := billing_routes.InitBillingAppRouter(conf)
	logger.Info("start auth app", zap.String("url", conf.AppHost+":"+conf.AppPort))
	if err := router.InitRoutes(conf.JWTKey).Start(":" + conf.AppPort); err != nil {
		logger.Fatal("Common server error", err)
	}
}
