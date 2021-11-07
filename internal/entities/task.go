package entities

import "github.com/google/uuid"

const (
	TaskCUDBrokerTopic = "taskCUD"
	TaskCreatedEvent   = "TaskCreated"
	NewTaskStatus      = "new"
)

type Task struct {
	TaskID      int64     `json:"-" db:"task_id"`
	PublicID    uuid.UUID `json:"public_id" db:"public_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Author      uuid.UUID `json:"author" db:"author"`
	AssignCost  int64     `json:"assign_cost" db:"assign_cost"`
	DoneCost    int64     `json:"done_cost" db:"done_cost"`
	Status      string    `json:"status" db:"status"`
}
