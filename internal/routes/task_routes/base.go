package task_routes

import (
	"net/http"

	"async_arch/internal/config"
	"async_arch/internal/service/task"

	"github.com/labstack/echo/v4"
)

type AuthAppRouter struct {
	httpEngine *echo.Echo
	appConf    *config.AppConfig
	uService   *task.UserTaskService
}

func InitAuthAppRouter(appConf *config.AppConfig, uService *task.UserTaskService) *AuthAppRouter {
	return &AuthAppRouter{
		httpEngine: echo.New(),
		appConf:    appConf,
		uService:   uService,
	}
}

func (ar *AuthAppRouter) InitRoutes(jwtKey string) *echo.Echo {
	ar.httpEngine.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	return ar.httpEngine
}
