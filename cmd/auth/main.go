package main

import (
	"flag"

	"async_arch/internal/config"
	"async_arch/internal/logger"
	"async_arch/internal/repository/user"
	"async_arch/internal/routes/auth_routes"
	"async_arch/internal/storage/database"

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
	router := auth_routes.InitAuthAppRouter(
		user.InitUserRepo(conn),
	)
	logger.Info("start auth app", zap.String("url", "http://localhost:"+conf.AppPort))
	if err := router.InitRoutes().Start(":" + conf.AppPort); err != nil {
		logger.Fatal("Common server error", err)
	}

}
