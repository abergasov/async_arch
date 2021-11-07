package task

import (
	"context"

	"async_arch/internal/logger"
	"async_arch/internal/repository/user"

	"github.com/segmentio/kafka-go"
)

type UserTaskService struct {
	uRepo  user.UserRepo
	jwtKey []byte
	broker *kafka.Reader
}

func InitUserTaskService(uRepo user.UserRepo, kfk *kafka.Reader, jwtKey string) *UserTaskService {
	u := &UserTaskService{uRepo: uRepo, jwtKey: []byte(jwtKey), broker: kfk}
	go u.readUserCUD()
	return u
}

func (u *UserTaskService) readUserCUD() {
	for {
		m, err := u.broker.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error read message from brocker", err)
			break
		}
		println(m.Key)
		println(m.Value)
	}
	if err := u.broker.Close(); err != nil {
		logger.Error("error broker closing", err)
	}
}
