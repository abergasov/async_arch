package service

import (
	"context"
	"errors"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/abergasov/schema_registry"
	"github.com/abergasov/schema_registry/pkg/grpc/user"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TaskUserRepo interface {
	CreateUser(account *user.UserAccountV1) error
	UpdateUser(account *user.UserAccountV1) error
}

type UserReplicatorService struct {
	uRepo    TaskUserRepo
	jwtKey   []byte
	broker   *kafka.Reader
	registry schema_registry.SchemaRegistry
}

func InitUserReplicatorService(uRepo TaskUserRepo, regio schema_registry.SchemaRegistry, kfk *kafka.Reader, jwtKey string) *UserReplicatorService {
	u := &UserReplicatorService{
		uRepo:    uRepo,
		jwtKey:   []byte(jwtKey),
		broker:   kfk,
		registry: regio,
	}
	go u.readUserCUD()
	return u
}

func (u *UserReplicatorService) readUserCUD() {
	for {
		m, err := u.broker.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error read message from broker", err)
			break
		}

		data, err := u.registry.DecodeUserStreamEvent(m.Value)
		if err != nil {
			logger.Error("error parse user from broker", err)
			continue
		}
		iUsr, ok := data[1]
		if !ok {
			logger.Error("error load user from map", errors.New("entity not in map"))
			continue
		}
		usr, ok := iUsr.(*user.UserAccountV1)
		if !ok {
			logger.Error("error cast user to user.UserAccountV1", errors.New("wrong type"))
			continue
		}
		switch string(m.Key) {
		case entities.UserCreatedEvent:
			u.createUser(usr)
		case entities.UserUpdatedEvent:
			u.updateUser(usr)
		}
	}
	if err := u.broker.Close(); err != nil {
		logger.Error("error broker closing", err)
	}
}

func (u *UserReplicatorService) createUser(usr *user.UserAccountV1) {
	if err := u.uRepo.CreateUser(usr); err != nil {
		logger.Error("error parse user from broker", err)
		return
	}
	logger.Info("new user created", zap.String("uuid", usr.PublicID))
}

func (u *UserReplicatorService) updateUser(usr *user.UserAccountV1) {
	if err := u.uRepo.UpdateUser(usr); err != nil {
		logger.Error("error parse user from broker", err)
		return
	}
	logger.Info("user updated", zap.String("uuid", usr.PublicID))
}
