package task

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/lib/pq"

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
	err := t.conn.Client().QueryRowx(`SELECT t.*, ta.assigned_to, ta.assigned_at 
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE public_id = $1`, taskID).StructScan(&tsk)
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
	rows, err := t.conn.Client().Queryx(`SELECT t.*, ta.assigned_to
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE DATE(t.created_at) = CURRENT_DATE`)
	return t.getTasks(rows, err)
}

func (t *TaskRepo) GetUserTasks(userPublicID uuid.UUID) ([]*entities.Task, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.*, ta.assigned_to
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE (ta.assigned_to = $1 OR t.author = $2) AND DATE(t.created_at) = CURRENT_DATE`, userPublicID, userPublicID)
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

func (t *TaskRepo) GetUnAssignedTasks() ([]*entities.Task, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.*, ta.assigned_to, ta.assigned_at
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE assigned_to IS NULL AND DATE(created_at) = CURRENT_DATE`)
	if err != nil {
		logger.Error("error load unassigned tasks", err)
		return nil, err
	}
	result := make([]*entities.Task, 0, 1000)
	defer rows.Close()
	for rows.Next() {
		var ts entities.Task
		if err = rows.StructScan(&ts); err != nil {
			logger.Error("error parse unassigned task", err)
			continue
		}
		result = append(result, &ts)
	}
	return result, nil
}

func (t *TaskRepo) AssignTasks(assign []*entities.TaskAssignContainer) error {
	sqlA := make([]string, 0, len(assign))
	sqlParams := make([]interface{}, 0, len(assign)*2)
	counter := 1

	taskUUIDs := make([]interface{}, 0, len(assign))
	for i := range assign {
		// placeholders for assigment
		sqlA = append(sqlA, fmt.Sprintf("($%d, $%d, NOW())", counter, counter+1))
		sqlParams = append(sqlParams, assign[i].TaskPublicID, assign[i].UserPublicID)
		counter += 2

		// placeholders for remove existing assigmants
		taskUUIDs = append(taskUUIDs, assign[i].TaskPublicID)
	}
	// clear old assignments to avoid duplicate errors
	if _, err := t.conn.Client().Exec(
		"DELETE FROM task_assignments WHERE task_uuid = any ($1)",
		pq.Array(taskUUIDs),
	); err != nil {
		logger.Error("error delete task assignments", err)
		return err
	}
	_, err := t.conn.Client().Exec("INSERT INTO task_assignments (task_uuid, assigned_to, assigned_at) VALUES "+strings.Join(sqlA, ","), sqlParams...)
	if err != nil {
		logger.Error("error task assign", err)
		return err
	}
	_, err = t.conn.Client().Exec(
		"UPDATE tasks SET status = $1 WHERE public_id = ANY ($2)",
		entities.AssignedTaskStatus,
		pq.Array(taskUUIDs),
	)
	if err != nil {
		logger.Error("error update task status", err)
	}
	return err
}

func (t *TaskRepo) DoneTask(taskPublicID uuid.UUID) error {
	_, err := t.conn.Client().Exec("UPDATE tasks SET status = $1 WHERE public_id = $2", entities.FinishTaskStatus, taskPublicID)
	return err
}
