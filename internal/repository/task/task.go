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

func (t *TaskRepo) GetByPublicID(taskID string) (*task.TaskV2, error) {
	var tsk task.TaskV2
	var assignedTo, assignedAt, trackerID, publicStatus sql.NullString
	err := t.conn.Client().QueryRowx(`SELECT t.tracker_id, t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.public_status,t.created_at, ta.assigned_to, ta.assigned_at 
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE public_id = $1`, taskID).
		Scan(&trackerID, &tsk.PublicID, &tsk.Author, &tsk.Title, &tsk.Description, &tsk.AssignCost, &tsk.DoneCost, &tsk.Status, &publicStatus, &tsk.CreatedAt, &assignedTo, &assignedAt)
	tsk.AssignedTo = assignedTo.String
	tsk.AssignedAt = assignedAt.String
	tsk.PublicStatus = publicStatus.String
	tsk.JiraID = trackerID.String
	if err != nil {
		logger.Error("error load task", err)
	}
	return &tsk, err
}

func (t *TaskRepo) CreateTask(taskAuthor, taskTitle, taskDesc, trackerID string) (*task.TaskV2, error) {
	newTaskID := uuid.New()
	assignCost, doneCost := t.calcCost()
	if _, err := t.conn.Client().NamedExec("INSERT INTO tasks (public_id,author,title,description,status,assign_cost,done_cost,created_at,tracker_id) VALUES (:public_id,:author,:title,:description,:status,:assign_cost,:done_cost,:created_at,:tracker_id)", map[string]interface{}{
		"public_id":   newTaskID,
		"author":      taskAuthor,
		"title":       taskTitle,
		"description": taskDesc,
		"status":      entities.NewTaskStatus,
		"assign_cost": assignCost,
		"created_at":  time.Now(),
		"done_cost":   doneCost,
		"tracker_id":  trackerID,
	}); err != nil {
		logger.Error("error task insert", err)
		return nil, err
	}
	return t.GetByPublicID(newTaskID.String())
}

func (t *TaskRepo) calcCost() (assignCost int64, doneCost int64) {
	return int64(rand.Intn(20-1) + 1), int64(rand.Intn(20-1) + 1)
}

func (t *TaskRepo) GetAllTasks() ([]*task.TaskV2, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.tracker_id,t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.public_status,t.created_at, ta.assigned_to
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE DATE(t.created_at) = CURRENT_DATE`)
	return t.getTasks(rows, err)
}

func (t *TaskRepo) GetUserTasks(userPublicID string) ([]*task.TaskV2, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.tracker_id,t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.public_status,t.created_at, ta.assigned_to
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE (ta.assigned_to = $1 OR t.author = $2) AND DATE(t.created_at) = CURRENT_DATE`, userPublicID, userPublicID)
	return t.getTasks(rows, err)
}

func (t *TaskRepo) getTasks(rows *sqlx.Rows, err error) ([]*task.TaskV2, error) {
	if err != nil {
		logger.Error("error load tasks", err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*task.TaskV2, 0, 100)
	for rows.Next() {
		var tsk task.TaskV2
		var assignedTo, trackerID, publicStatus sql.NullString
		if err = rows.Scan(&trackerID, &tsk.PublicID, &tsk.Author, &tsk.Title, &tsk.Description, &tsk.AssignCost, &tsk.DoneCost, &tsk.Status, &publicStatus, &tsk.CreatedAt, &assignedTo); err != nil {
			logger.Error("error scan task", err)
			continue
		}
		tsk.AssignedTo = assignedTo.String
		tsk.JiraID = trackerID.String
		tsk.PublicStatus = publicStatus.String
		result = append(result, &tsk)
	}
	return result, err
}

func (t *TaskRepo) GetUnAssignedTasks() ([]*task.TaskV2, error) {
	rows, err := t.conn.Client().Queryx(`SELECT t.tracker_id,t.public_id,t.author,t.title,t.description,t.assign_cost,t.done_cost,t.status,t.created_at, ta.assigned_to, ta.assigned_at
		FROM tasks t
		LEFT JOIN task_assignments ta ON ta.task_uuid = t.public_id
		WHERE assigned_to IS NULL AND DATE(created_at) = CURRENT_DATE`)
	if err != nil {
		logger.Error("error load unassigned tasks", err)
		return nil, err
	}
	result := make([]*task.TaskV2, 0, 1000)
	defer rows.Close()
	var assignedTo, assignedAt, trackerID sql.NullString
	for rows.Next() {
		var tsk task.TaskV2

		if err = rows.Scan(&trackerID, &tsk.PublicID, &tsk.Author, &tsk.Title, &tsk.Description, &tsk.AssignCost, &tsk.DoneCost, &tsk.Status, &tsk.CreatedAt, &assignedTo, &assignedAt); err != nil {
			logger.Error("error parse unassigned task", err)
			continue
		}
		tsk.AssignedTo = assignedTo.String
		tsk.AssignedAt = assignedAt.String
		tsk.JiraID = trackerID.String
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
		"UPDATE tasks SET status = $1, public_status = $2 WHERE public_id = ANY ($3)",
		entities.AssignedTaskStatus,
		entities.PublicAssignedTaskStatus,
		pq.Array(taskUUIDs),
	)
	if err != nil {
		logger.Error("error update task status", err)
	}
	return err
}

func (t *TaskRepo) FinishTask(taskPublicID string) error {
	_, err := t.conn.Client().Exec(
		"UPDATE tasks SET status = $1, public_status = $2 WHERE public_id = $3",
		entities.FinishTaskStatus,
		entities.PublicFinishTaskStatus,
		taskPublicID,
	)
	return err
}

func (t *TaskRepo) CreateTaskV1(tsk *task.TaskV1) error {
	_, err := t.conn.Client().NamedExec("INSERT INTO tasks (public_id, author, title, description, assign_cost, done_cost, status, created_at) VALUES (:public_id, :author, :title, :description, :assign_cost, :done_cost, :status, :created_at)", map[string]interface{}{
		"public_id":   tsk.PublicID,
		"author":      tsk.Author,
		"title":       tsk.Title,
		"description": tsk.Description,
		"assign_cost": tsk.AssignCost,
		"done_cost":   tsk.DoneCost,
		"status":      tsk.Status,
		"created_at":  tsk.CreatedAt,
	})
	return err
}

func (t *TaskRepo) CreateTaskV2(tsk *task.TaskV2) error {
	_, err := t.conn.Client().NamedExec("INSERT INTO tasks (public_id, author, title, description, assign_cost, done_cost, status, created_at, tracker_id, public_status) VALUES (:public_id, :author, :title, :description, :assign_cost, :done_cost, :status, :created_at, :tracker_id, :public_status)", map[string]interface{}{
		"public_id":     tsk.PublicID,
		"author":        tsk.Author,
		"title":         tsk.Title,
		"description":   tsk.Description,
		"assign_cost":   tsk.AssignCost,
		"done_cost":     tsk.DoneCost,
		"status":        tsk.Status,
		"created_at":    tsk.CreatedAt,
		"tracker_id":    tsk.JiraID,
		"public_status": tsk.PublicStatus,
	})
	return err
}

func (t *TaskRepo) UpdateTaskV1(tsk *task.TaskV1) error {
	_, err := t.conn.Client().NamedExec("UPDATE tasks SET author=:author, title=:title, description=:description, assign_cost=:assign_cost, done_cost=:done_cost, status=:status, created_at=:created_at WHERE public_id=:public_id", map[string]interface{}{
		"public_id":   tsk.PublicID,
		"author":      tsk.Author,
		"title":       tsk.Title,
		"description": tsk.Description,
		"assign_cost": tsk.AssignCost,
		"done_cost":   tsk.DoneCost,
		"status":      tsk.Status,
		"created_at":  tsk.CreatedAt,
	})
	return err
}

func (t *TaskRepo) UpdateTaskV2(tsk *task.TaskV2) error {
	_, err := t.conn.Client().NamedExec("UPDATE tasks SET author=:author, title=:title, description=:description, assign_cost=:assign_cost, done_cost=:done_cost, status=:status, created_at=:created_at, tracker_id:tracker_id, public_status:public_status WHERE public_id=:public_id", map[string]interface{}{
		"public_id":     tsk.PublicID,
		"author":        tsk.Author,
		"title":         tsk.Title,
		"description":   tsk.Description,
		"assign_cost":   tsk.AssignCost,
		"done_cost":     tsk.DoneCost,
		"status":        tsk.Status,
		"created_at":    tsk.CreatedAt,
		"tracker_id":    tsk.JiraID,
		"public_status": tsk.PublicStatus,
	})
	return err
}
