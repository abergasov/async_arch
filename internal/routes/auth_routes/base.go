package auth_routes

import (
	"net/http"

	"async_arch/internal/config"
	"async_arch/internal/repository/exchanger"
	"async_arch/internal/service"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
}

type AuthAppRouter struct {
	httpEngine  *echo.Echo
	googleOAuth *oauth2.Config
	appConf     *config.AppConfig
	uService    *service.UserService
	exchanger   *exchanger.Exchanger
}

func InitAuthAppRouter(appConf *config.AppConfig, uService *service.UserService, exchanger *exchanger.Exchanger) *AuthAppRouter {
	return &AuthAppRouter{
		httpEngine: echo.New(),
		exchanger:  exchanger,
		appConf:    appConf,
		googleOAuth: &oauth2.Config{
			RedirectURL:  appConf.AppHost + ":" + appConf.AppPort + "/auth/google/callback",
			ClientID:     appConf.GoogleAppID,
			ClientSecret: appConf.GoogleAppSecret,
			Scopes:       googleScopes,
			Endpoint:     google.Endpoint,
		},
		uService: uService,
	}
}

func (ar *AuthAppRouter) InitRoutes() *echo.Echo {
	ar.httpEngine.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	ar.httpEngine.GET("/auth/google/login", ar.oauthGoogleLogin)
	ar.httpEngine.GET("/auth/google/callback", ar.oauthGoogleCallback)
	ar.httpEngine.POST("/api/v1/exchange", ar.exchangeCode)
	return ar.httpEngine
}
