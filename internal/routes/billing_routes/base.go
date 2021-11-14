package billing_routes

import (
	"net/http"

	"async_arch/internal/config"

	"github.com/labstack/echo/v4"
)

type BillingAppRouter struct {
	httpEngine *echo.Echo
	appConf    *config.AppConfig
}

func InitBillingAppRouter(appConf *config.AppConfig) *BillingAppRouter {
	return &BillingAppRouter{
		httpEngine: echo.New(),
		appConf:    appConf,
	}
}

func (ar *BillingAppRouter) InitRoutes(jwtKey string) *echo.Echo {
	ar.httpEngine.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	return ar.httpEngine
}
