package service

import (
	"context"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/segmentio/kafka-go"
)

type TaskReplicatorService struct {
	uRepo  TaskUserRepo
	jwtKey []byte
	broker *kafka.Reader
}

func InitTaskReplicatorService(uRepo TaskUserRepo, kfk *kafka.Reader, jwtKey string) *TaskReplicatorService {
	t := &TaskReplicatorService{uRepo: uRepo, jwtKey: []byte(jwtKey), broker: kfk}
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

		switch string(m.Key) {
		case entities.TaskCreatedEvent:
			t.createTask(m.Value)
		case entities.TaskAssignedEvent:
			t.assignTask(m.Value)
		case entities.TaskFinishEvent:
			t.finishTask(m.Value)
		}
	}
	if err := t.broker.Close(); err != nil {
		logger.Error("error broker closing", err)
	}
}

func (t *TaskReplicatorService) finishTask(rawData []byte) {}

func (t *TaskReplicatorService) assignTask(rawData []byte) {}

func (t *TaskReplicatorService) createTask(rawData []byte) {}
