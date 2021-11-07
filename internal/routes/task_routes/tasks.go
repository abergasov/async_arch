package task_routes

import (
	"net/http"

	"async_arch/internal/entities"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (ar *TaskAppRouter) createTask(c echo.Context) error {
	var t struct {
		Task struct {
			Title string `json:"title"`
			Desc  string `json:"desc"`
		} `json:"task"`
	}
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, entities.ErrorRequest{})
	}
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	newTask, err := ar.tManager.CreateTask(claims.UserID, t.Task.Title, t.Task.Title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entities.ErrorRequest{})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true, "data": newTask})
}
