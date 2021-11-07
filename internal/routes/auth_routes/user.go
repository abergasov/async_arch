package auth_routes

import (
	"net/http"

	"async_arch/internal/entities"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (ar *AuthAppRouter) changeRole(c echo.Context) error {
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	var u struct {
		Role string `json:"role"`
	}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]bool{"ok": false})
	}
	usr, jwtKey, err := ar.uService.ChangeRole(claims.UserID, claims.UserVersion, u.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]bool{"ok": false})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true, "data": usr, "token": jwtKey})
}
