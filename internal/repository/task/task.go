package task

import (
	"math/rand"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/google/uuid"
)

type TaskRepo struct {
	conn database.DBConnector
}

func InitTaskRepo(conn database.DBConnector) *TaskRepo {
	return &TaskRepo{conn: conn}
}

func (t *TaskRepo) GetTaskByPublicID(taskID uuid.UUID) (*entities.Task, error) {
	sqlS := "SELECT task_id,public_id,author,title,description,assign_cost,status,done_cost FROM tasks WHERE public_id = $1"
	var tsk entities.Task
	err := t.conn.Client().QueryRowx(sqlS, taskID).StructScan(&tsk)
	return &tsk, err
}

func (t *TaskRepo) CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error) {
	newTaskID := uuid.New()
	assignCost, doneCost := t.calcCost()
	if _, err := t.conn.Client().NamedExec("INSERT INTO tasks (public_id,author,title,description,status,assign_cost,done_cost) VALUES (:public_id,:author,:title,:description,:status,:assign_cost,:done_cost)", map[string]interface{}{
		"public_id":   newTaskID,
		"author":      taskAuthor,
		"title":       taskTitle,
		"description": taskDesc,
		"status":      entities.NewTaskStatus,
		"assign_cost": assignCost,
		"done_cost":   doneCost,
	}); err != nil {
		logger.Error("error task insert", err)
		return nil, err
	}
	return t.GetTaskByPublicID(newTaskID)
}

func (t *TaskRepo) calcCost() (assignCost int64, doneCost int64) {
	return int64(rand.Intn(20-1) + 1), int64(rand.Intn(20-1) + 1)
}
