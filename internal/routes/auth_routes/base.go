package auth_routes

import (
	"net/http"

	"async_arch/internal/repository/user"

	"github.com/labstack/echo/v4"
)

type AuthAppRouter struct {
	httpEngine *echo.Echo
	userRepo   user.UserRepo
}

func InitAuthAppRouter(uR user.UserRepo) *AuthAppRouter {
	return &AuthAppRouter{
		httpEngine: echo.New(),
		userRepo:   uR,
	}
}

func (ar *AuthAppRouter) InitRoutes() *echo.Echo {
	ar.httpEngine.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	return ar.httpEngine
}
