package task

import (
	"context"
	"encoding/json"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/service/auth"

	"github.com/segmentio/kafka-go"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserByPublicID(publicID uuid.UUID, version int) (*entities.UserAccount, error)
}

type TaskRepo interface {
	GetTaskByPublicID(taskID uuid.UUID) (*entities.Task, error)
	CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error)
	GetAllTasks() ([]*entities.Task, error)
	GetUserTasks(userPublicID uuid.UUID) ([]*entities.Task, error)
}

type TaskManager struct {
	tRepo  TaskRepo
	uRepo  UserRepo
	broker *kafka.Writer
}

func InitTaskManager(t TaskRepo, u UserRepo, kfk *kafka.Writer) *TaskManager {
	return &TaskManager{
		tRepo:  t,
		uRepo:  u,
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

func (t *TaskManager) LoadTasks(userPublicID uuid.UUID, userVersion int) ([]*entities.Task, error) {
	usr, err := t.uRepo.GetUserByPublicID(userPublicID, userVersion)
	if err != nil {
		return nil, err
	}
	if usr.UserRole == auth.Worker {
		return t.tRepo.GetUserTasks(userPublicID)
	}
	return t.tRepo.GetAllTasks()
}
