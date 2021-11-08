package task

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/service/auth"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type UserRepo interface {
	GetActiveWorkers() ([]uuid.UUID, error)
	GetUserByPublicID(publicID uuid.UUID, version int) (*entities.UserAccount, error)
}

type TaskRepo interface {
	GetTaskByPublicID(taskID uuid.UUID) (*entities.Task, error)
	CreateTask(taskAuthor uuid.UUID, taskTitle, taskDesc string) (*entities.Task, error)
	GetAllTasks() ([]*entities.Task, error)
	GetUserTasks(userPublicID uuid.UUID) ([]*entities.Task, error)
	GetUnAssignedTasks() ([]*entities.Task, error)
	AssignTasks(assign []*entities.TaskAssignContainer) error
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

func (t *TaskManager) AssignTasks(userPublicID uuid.UUID, userVersion int) ([]*entities.Task, error) {
	usr, err := t.uRepo.GetUserByPublicID(userPublicID, userVersion)
	if err != nil {
		return nil, err
	}
	if usr.UserRole == auth.Worker {
		err = errors.New("only admin and manager can assign tasks")
		logger.Error("error assign tasks", err)
		return nil, err
	}
	workers, err := t.uRepo.GetActiveWorkers()
	if err != nil {
		logger.Error("error get active workers", err)
		return nil, err
	}
	tasks, err := t.tRepo.GetUnAssignedTasks()
	if err != nil {
		logger.Error("error get unassigned tasks", err)
		return nil, err
	}
	if len(tasks) == 0 {
		return t.LoadTasks(userPublicID, userVersion)
	}
	if len(workers) == 0 {
		err = errors.New("empty workers list")
		logger.Error("there is no workers in system", err)
		return nil, err
	}
	if err = t.assignTasks(workers, tasks); err != nil {
		logger.Error("error assign tasks", err)
		return nil, err
	}
	return t.LoadTasks(userPublicID, userVersion)
}

func (t *TaskManager) assignTasks(userIDs []uuid.UUID, targetTasks []*entities.Task) error {
	assigned := make([]*entities.TaskAssignContainer, 0, len(targetTasks))
	for i := range targetTasks {
		assigned = append(assigned, &entities.TaskAssignContainer{
			TaskPublicID: targetTasks[i].PublicID,
			UserPublicID: userIDs[rand.Intn(len(userIDs)-0)+0],
		})
	}
	return t.tRepo.AssignTasks(assigned)
}
