package task

import (
	"math/rand"
	"time"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
)

type TaskRepo struct {
	conn database.DBConnector
}

func InitTaskRepo(conn database.DBConnector) *TaskRepo {
	return &TaskRepo{conn: conn}
}

func (t *TaskRepo) GetTaskByPublicID(taskID uuid.UUID) (*entities.Task, error) {
	var tsk entities.Task
	err := t.conn.Client().QueryRowx("SELECT * FROM tasks WHERE public_id = $1", taskID).StructScan(&tsk)
	if err != nil {
		logger.Error("error load task", err)
	}
	return &tsk, err
}

func (t *TaskRepo) CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error) {
	newTaskID := uuid.New()
	assignCost, doneCost := t.calcCost()
	if _, err := t.conn.Client().NamedExec("INSERT INTO tasks (public_id,author,title,description,status,assign_cost,done_cost,created_at) VALUES (:public_id,:author,:title,:description,:status,:assign_cost,:done_cost,:created_at)", map[string]interface{}{
		"public_id":   newTaskID,
		"author":      taskAuthor,
		"title":       taskTitle,
		"description": taskDesc,
		"status":      entities.NewTaskStatus,
		"assign_cost": assignCost,
		"created_at":  time.Now(),
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

func (t *TaskRepo) GetAllTasks() ([]*entities.Task, error) {
	rows, err := t.conn.Client().Queryx("SELECT * FROM tasks WHERE DATE(created_at) = CURRENT_DATE")
	return t.getTasks(rows, err)
}

func (t *TaskRepo) GetUserTasks(userPublicID uuid.UUID) ([]*entities.Task, error) {
	rows, err := t.conn.Client().Queryx("SELECT * FROM tasks WHERE (assigned_to = $1 OR author = $2) AND DATE(created_at) = CURRENT_DATE", userPublicID, userPublicID)
	return t.getTasks(rows, err)
}

func (t *TaskRepo) getTasks(rows *sqlx.Rows, err error) ([]*entities.Task, error) {
	if err != nil {
		logger.Error("error load tasks", err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*entities.Task, 0, 100)
	for rows.Next() {
		var tsk entities.Task
		if err = rows.StructScan(&tsk); err != nil {
			logger.Error("error scan task", err)
			continue
		}
		result = append(result, &tsk)
	}
	return result, err
}
