package task

import (
	"context"
	"encoding/json"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TaskUserRepo interface {
	CreateUser(account *entities.UserAccount) error
	UpdateUser(account *entities.UserAccount) error
}

type UserTaskService struct {
	uRepo  TaskUserRepo
	jwtKey []byte
	broker *kafka.Reader
}

func InitUserTaskService(uRepo TaskUserRepo, kfk *kafka.Reader, jwtKey string) *UserTaskService {
	u := &UserTaskService{uRepo: uRepo, jwtKey: []byte(jwtKey), broker: kfk}
	go u.readUserCUD()
	return u
}

func (u *UserTaskService) readUserCUD() {
	for {
		m, err := u.broker.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error read message from broker", err)
			break
		}
		switch string(m.Key) {
		case entities.UserCreatedEvent:
			u.createUser(m.Value)
		case entities.UserUpdatedEvent:
			u.updateUser(m.Value)
		}
	}
	if err := u.broker.Close(); err != nil {
		logger.Error("error broker closing", err)
	}
}

func (u *UserTaskService) createUser(rawData []byte) {
	var usr entities.UserAccount
	if err := json.Unmarshal(rawData, &usr); err != nil {
		logger.Error("error parse user from broker", err)
		return
	}
	if err := u.uRepo.CreateUser(&usr); err != nil {
		logger.Error("error parse user from broker", err)
		return
	}
	logger.Info("new user created", zap.String("uuid", usr.PublicID.String()))
}

func (u *UserTaskService) updateUser(rawData []byte) {
	var usr entities.UserAccount
	if err := json.Unmarshal(rawData, &usr); err != nil {
		logger.Error("error parse user from broker", err)
		return
	}
	if err := u.uRepo.UpdateUser(&usr); err != nil {
		logger.Error("error parse user from broker", err)
		return
	}
	logger.Info("user updated", zap.String("uuid", usr.PublicID.String()))
}
