package task

import (
	"context"
	"encoding/json"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/segmentio/kafka-go"

	"github.com/google/uuid"
)

type TaskRepo interface {
	GetTaskByPublicID(taskID uuid.UUID) (*entities.Task, error)
	CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error)
}

type TaskManager struct {
	tRepo  TaskRepo
	broker *kafka.Writer
}

func InitTaskManager(t TaskRepo, kfk *kafka.Writer) *TaskManager {
	return &TaskManager{
		tRepo:  t,
		broker: kfk,
	}
}

func (t *TaskManager) CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error) {
	task, err := t.tRepo.CreateTask(taskAuthor, taskTitle, taskDesc)
	if err != nil {
		logger.Error("error create task", err)
		return nil, err
	}
	b, _ := json.Marshal(task)
	if err = t.broker.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(entities.TaskCreatedEvent),
		Value: b,
	}); err != nil {
		logger.Error("error stream task", err)
		return nil, err
	}
	return task, nil
}
