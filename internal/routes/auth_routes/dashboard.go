package auth_routes

import (
	"net/http"

	"async_arch/internal/entities"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (ar *AuthAppRouter) dashboardData(c echo.Context) error {
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	usr, err := ar.uService.GetUserInfo(claims.UserID, claims.UserVersion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]bool{"ok": true})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
		"data": map[string]interface{}{
			"name":      usr.UserName,
			"email":     usr.UserMail,
			"public_id": usr.PublicID,
			"role":      usr.UserRole,
		},
	})
}
