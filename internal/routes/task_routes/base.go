package task_routes

import (
	"net/http"

	"async_arch/internal/config"
	"async_arch/internal/entities"
	"async_arch/internal/service/task"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

type TaskAppRouter struct {
	httpEngine *echo.Echo
	appConf    *config.AppConfig
	tManager   *task.TaskManager
}

func InitAuthAppRouter(appConf *config.AppConfig, tManager *task.TaskManager) *TaskAppRouter {
	return &TaskAppRouter{
		httpEngine: echo.New(),
		appConf:    appConf,
		tManager:   tManager,
	}
}

func (ar *TaskAppRouter) InitRoutes(jwtKey string) *echo.Echo {
	ar.httpEngine.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	userData := ar.httpEngine.Group("/api/v1/task/")
	{
		userData.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:        &entities.UserJWT{},
			SigningKey:    []byte(jwtKey),
			SigningMethod: jwt.SigningMethodHS512.Name,
		}))
		userData.POST("create", ar.createTask)
		userData.POST("list", ar.getTaskList)
		userData.POST("assign", ar.assignFreeTasks)
		userData.POST("finish", ar.finish)
	}
	return ar.httpEngine
}
