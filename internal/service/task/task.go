package task

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/abergasov/schema_registry/pkg/grpc/task"
	"github.com/abergasov/schema_registry/pkg/grpc/user"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/service/auth"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type UserRepo interface {
	GetActiveWorkers() ([]uuid.UUID, error)
	GetByPublicID(publicID string, version int64) (*user.UserAccountV1, error)
}

type TaskRepo interface {
	GetByPublicID(taskID string) (*task.TaskV1, error)
	CreateTask(taskAuthor, taskTitle, taskDesc string) (*task.TaskV1, error)
	GetAllTasks() ([]*task.TaskV1, error)
	GetUserTasks(userPublicID string) ([]*task.TaskV1, error)
	GetUnAssignedTasks() ([]*task.TaskV1, error)
	AssignTasks(assign []*entities.TaskAssignContainer) error
	FinishTask(taskPublicID string) error
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

func (t *TaskManager) CreateTask(taskAuthor string, taskTitle, taskDesc string) (*task.TaskV1, error) {
	tsk, err := t.tRepo.CreateTask(taskAuthor, taskTitle, taskDesc)
	if err != nil {
		logger.Error("error create task", err)
		return nil, err
	}
	b, _ := json.Marshal(tsk)
	if err = t.broker.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(entities.TaskCreatedEvent),
		Value: b,
	}); err != nil {
		logger.Error("error stream task", err)
		return nil, err
	}
	return tsk, nil
}

func (t *TaskManager) LoadTasks(userPublicID string, userVersion int64) ([]*task.TaskV1, error) {
	usr, err := t.uRepo.GetByPublicID(userPublicID, userVersion)
	if err != nil {
		return nil, err
	}
	if usr.UserRole == auth.Worker {
		return t.tRepo.GetUserTasks(userPublicID)
	}
	return t.tRepo.GetAllTasks()
}

func (t *TaskManager) AssignTasks(userPublicID string, userVersion int64) ([]*task.TaskV1, error) {
	usr, err := t.uRepo.GetByPublicID(userPublicID, userVersion)
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
	messages := make([]kafka.Message, 0, len(tasks))
	for i := range tasks {
		b, _ := json.Marshal(tasks[i])
		messages = append(messages, kafka.Message{
			Key:   []byte(entities.TaskCreatedEvent),
			Value: b,
		})
	}
	if err = t.broker.WriteMessages(context.Background(), messages...); err != nil {
		logger.Error("error stream event", err)
		return nil, err
	}

	return t.LoadTasks(userPublicID, userVersion)
}

func (t *TaskManager) assignTasks(userIDs []uuid.UUID, targetTasks []*task.TaskV1) error {
	assigned := make([]*entities.TaskAssignContainer, 0, len(targetTasks))
	for i := range targetTasks {
		workerID := userIDs[rand.Intn(len(userIDs)-0)+0]
		assigned = append(assigned, &entities.TaskAssignContainer{
			TaskPublicID: targetTasks[i].PublicID,
			UserPublicID: workerID.String(),
		})
		targetTasks[i].AssignedAt = time.Now().Format(time.RFC3339)
		targetTasks[i].AssignedTo = workerID.String()
	}
	return t.tRepo.AssignTasks(assigned)
}

func (t *TaskManager) Finish(taskPublicID, userPublicID string, userVersion int64) error {
	usr, err := t.uRepo.GetByPublicID(userPublicID, userVersion)
	if err != nil {
		return err
	}
	targetTask, err := t.tRepo.GetByPublicID(taskPublicID)
	if err != nil {
		logger.Error("error load task by publicID", err)
		return err
	}
	if targetTask.AssignedTo != usr.PublicID {
		err = errors.New("task not assigned to this user")
		return err
	}
	if err = t.tRepo.FinishTask(taskPublicID); err != nil {
		logger.Error("erorr finish task", err)
		return err
	}
	b, _ := json.Marshal(targetTask)
	if err = t.broker.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(entities.TaskFinishEvent),
		Value: b,
	}); err != nil {
		logger.Error("error stream event", err)
		return err
	}
	return nil
}
