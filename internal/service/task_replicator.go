package service

import (
	"context"
	"errors"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/abergasov/schema_registry"
	"github.com/abergasov/schema_registry/pkg/grpc/task"

	"github.com/segmentio/kafka-go"
)

type TaskRepo interface {
	CreateTaskV1(*task.TaskV1) error
	CreateTaskV2(*task.TaskV2) error
	UpdateTaskV1(*task.TaskV1) error
	UpdateTaskV2(*task.TaskV2) error
}

type TaskReplicatorService struct {
	tRepo    TaskRepo
	jwtKey   []byte
	broker   *kafka.Reader
	brokerBI *kafka.Reader
	registry schema_registry.SchemaRegistry
}

func InitTaskReplicatorService(tRepo TaskRepo, regio schema_registry.SchemaRegistry, kfk *kafka.Reader, kfkBI *kafka.Reader) *TaskReplicatorService {
	t := &TaskReplicatorService{
		tRepo:    tRepo,
		broker:   kfk,
		brokerBI: kfkBI,
		registry: regio,
	}
	go t.readTaskCUD()
	return t
}

func (t *TaskReplicatorService) readTaskCUD() {
	for {
		m, err := t.broker.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error read message from broker", err)
			break
		}

		data, err := t.registry.DecodeTaskStreamEvent(m.Value)
		if err != nil {
			logger.Error("error parse task from broker", err)
			continue
		}

		if tObj, ok := data[1]; ok {
			tV1, ok := tObj.(*task.TaskV1)
			if !ok {
				logger.Error("error cast message to taskV1", errors.New("wrong type"))
				continue
			}
			t.processV1(string(m.Key), tV1)
		}

		if tObj, ok := data[2]; ok {
			tV2, ok := tObj.(*task.TaskV2)
			if !ok {
				logger.Error("error cast message to taskV1", errors.New("wrong type"))
				continue
			}
			t.processV2(string(m.Key), tV2)
		}
	}
	if err := t.broker.Close(); err != nil {
		logger.Error("error broker closing", err)
	}
}

func (t *TaskReplicatorService) processV1(event string, task *task.TaskV1) {
	var err error
	switch event {
	case entities.TaskCreatedEvent:
		err = t.tRepo.CreateTaskV1(task)
	case entities.TaskAssignedEvent, entities.TaskFinishEvent:
		err = t.tRepo.UpdateTaskV1(task)
	}
	if err != nil {
		logger.Error("error process taskV1", err)
	}
}

func (t *TaskReplicatorService) processV2(event string, task *task.TaskV2) {
	var err error
	switch event {
	case entities.TaskCreatedEvent:
		err = t.tRepo.CreateTaskV2(task)
	case entities.TaskAssignedEvent, entities.TaskFinishEvent:
		err = t.tRepo.UpdateTaskV2(task)
	}
	if err != nil {
		logger.Error("error process taskV2", err)
	}
}
