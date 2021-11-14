package task

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/abergasov/schema_registry/pkg/grpc/task"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TaskRepo struct {
	conn database.DBConnector
}

func InitTaskRepo(conn database.DBConnector) *TaskRepo {
	return &TaskRepo{conn: conn}
}

func (t *TaskRepo) GetByPublicID(taskID string) (*task.TaskV1, error) {
	var tsk task.TaskV1
	var assignedTo, assignedAt sql.NullString
	err := t.conn.Client().QueryRowx(`SELECT t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.created_at, ta.assigned_to, ta.assigned_at 
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE public_id = $1`, taskID).
		Scan(&tsk.PublicID, &tsk.Author, &tsk.Title, &tsk.Description, &tsk.AssignCost, &tsk.DoneCost, &tsk.Status, &tsk.CreatedAt, &assignedTo, &assignedAt)
	tsk.AssignedTo = assignedTo.String
	tsk.AssignedAt = assignedAt.String
	if err != nil {
		logger.Error("error load task", err)
	}
	return &tsk, err
}

func (t *TaskRepo) CreateTask(taskAuthor, taskTitle, taskDesc string) (*task.TaskV1, error) {
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
	return t.GetByPublicID(newTaskID.String())
}

func (t *TaskRepo) calcCost() (assignCost int64, doneCost int64) {
	return int64(rand.Intn(20-1) + 1), int64(rand.Intn(20-1) + 1)
}

func (t *TaskRepo) GetAllTasks() ([]*task.TaskV1, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.created_at, ta.assigned_to
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE DATE(t.created_at) = CURRENT_DATE`)
	return t.getTasks(rows, err)
}

func (t *TaskRepo) GetUserTasks(userPublicID string) ([]*task.TaskV1, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.created_at, ta.assigned_to
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE (ta.assigned_to = $1 OR t.author = $2) AND DATE(t.created_at) = CURRENT_DATE`, userPublicID, userPublicID)
	return t.getTasks(rows, err)
}

func (t *TaskRepo) getTasks(rows *sqlx.Rows, err error) ([]*task.TaskV1, error) {
	if err != nil {
		logger.Error("error load tasks", err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*task.TaskV1, 0, 100)
	for rows.Next() {
		var tsk task.TaskV1
		var assignedTo sql.NullString
		if err = rows.Scan(&tsk.PublicID, &tsk.Author, &tsk.Title, &tsk.Description, &tsk.AssignCost, &tsk.DoneCost, &tsk.Status, &tsk.CreatedAt, &assignedTo); err != nil {
			logger.Error("error scan task", err)
			continue
		}
		tsk.AssignedTo = assignedTo.String
		result = append(result, &tsk)
	}
	return result, err
}

func (t *TaskRepo) GetUnAssignedTasks() ([]*task.TaskV1, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.created_at, ta.assigned_to, ta.assigned_at
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE assigned_to IS NULL AND DATE(created_at) = CURRENT_DATE`)
	if err != nil {
		logger.Error("error load unassigned tasks", err)
		return nil, err
	}
	result := make([]*task.TaskV1, 0, 1000)
	defer rows.Close()
	var assignedTo, assignedAt sql.NullString
	for rows.Next() {
		var tsk task.TaskV1

		if err = rows.Scan(&tsk.PublicID, &tsk.Author, &tsk.Title, &tsk.Description, &tsk.AssignCost, &tsk.DoneCost, &tsk.Status, &tsk.CreatedAt, &assignedTo, &assignedAt); err != nil {
			logger.Error("error parse unassigned task", err)
			continue
		}
		tsk.AssignedTo = assignedTo.String
		tsk.AssignedAt = assignedAt.String
		result = append(result, &tsk)
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

func (t *TaskRepo) FinishTask(taskPublicID string) error {
	_, err := t.conn.Client().Exec("UPDATE tasks SET status = $1 WHERE public_id = $2", entities.FinishTaskStatus, taskPublicID)
	return err
}
