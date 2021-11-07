package task

import (
	"async_arch/internal/entities"

	"github.com/google/uuid"
)

type TaskRepo interface {
	GetTaskByPublicID(taskID uuid.UUID) (*entities.Task, error)
	CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error)
}

type TaskManager struct {
	tRepo TaskRepo
}

func InitTaskManager(t TaskRepo) *TaskManager {
	return &TaskManager{
		tRepo: t,
	}
}

func (t *TaskManager) CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error) {
	return t.tRepo.CreateTask(taskAuthor, taskTitle, taskDesc)
}
