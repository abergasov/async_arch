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
	if t.Task.Title == "" || t.Task.Desc == "" {
		return c.JSON(http.StatusBadRequest, entities.ErrorRequest{})
	}
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	newTask, err := ar.tManager.CreateTask(claims.UserID, t.Task.Title, t.Task.Desc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entities.ErrorRequest{})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true, "data": newTask})
}

func (ar *TaskAppRouter) getTaskList(c echo.Context) error {
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	tasks, err := ar.tManager.LoadTasks(claims.UserID, claims.UserVersion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entities.ErrorRequest{})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true, "data": tasks})
}

func (ar *TaskAppRouter) assignFreeTasks(c echo.Context) error {
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	tasks, err := ar.tManager.AssignTasks(claims.UserID, claims.UserVersion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entities.ErrorRequest{})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true, "data": tasks})
}

func (ar *TaskAppRouter) finish(c echo.Context) error {
	user, _ := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(*entities.UserJWT)
	var t struct {
		TaskID string `json:"task_id"`
	}
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, entities.ErrorRequest{})
	}
	err := ar.tManager.Finish(t.TaskID, claims.UserID, claims.UserVersion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entities.ErrorRequest{})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true})
}
