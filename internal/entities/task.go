package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const (
	TaskCUDBrokerTopic = "taskStream"
	TaskCreatedEvent   = "TaskCreated"
	TaskAssignedEvent  = "TaskAssigned"
	TaskFinishEvent    = "TaskFinished"
	NewTaskStatus      = "new"
	AssignedTaskStatus = "assigned"
	FinishTaskStatus   = "done"
)

type Task struct {
	TaskID      int64        `json:"-" db:"task_id"`
	PublicID    uuid.UUID    `json:"public_id" db:"public_id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Author      uuid.UUID    `json:"author" db:"author"`
	AssignCost  int64        `json:"assign_cost" db:"assign_cost"`
	DoneCost    int64        `json:"done_cost" db:"done_cost"`
	Status      string       `json:"status" db:"status"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	AssignedTo  uuid.UUID    `json:"assigned_to" db:"assigned_to"`
	AssignedAt  sql.NullTime `json:"assigned_at" db:"assigned_at"`
}

type TaskAssignContainer struct {
	TaskPublicID uuid.UUID
	UserPublicID uuid.UUID
}
