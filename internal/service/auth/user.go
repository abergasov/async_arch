package auth

import (
	"context"
	"encoding/json"
	"time"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	userRepo "async_arch/internal/repository/user"

	"github.com/abergasov/schema_registry"

	"github.com/abergasov/schema_registry/pkg/grpc/user"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
)

const (
	Worker    = "worker"
	AdminRole = "admin"
	Manager   = "manager"
)

type UserService struct {
	uRepo    userRepo.UserRepo
	jwtKey   []byte
	broker   *kafka.Writer
	brokerBI *kafka.Writer
	registry schema_registry.SchemaRegistry
}

func InitUserService(
	uRepo userRepo.UserRepo,
	regio schema_registry.SchemaRegistry,
	kfk *kafka.Writer,
	kfkBI *kafka.Writer,
	jwtKey string,
) *UserService {
	return &UserService{uRepo: uRepo, jwtKey: []byte(jwtKey), broker: kfk, brokerBI: kfkBI, registry: regio}
}

func (u *UserService) Login(googleUser *entities.GoogleUser) (string, error) {
	usr, err := u.uRepo.GetUserByMail(googleUser.Email)
	if err != nil {
		logger.Error("error load user by mail", err)
		return "", err
	}
	if usr == nil {
		if err = u.registerUser(googleUser); err != nil {
			return "", err
		}
		usr, err = u.uRepo.GetUserByMail(googleUser.Email)
		if err != nil {
			logger.Error("error load user by mail after creation", err)
			return "", err
		}

		// stream create event to broker
		b, err := u.registry.EncodeUserStreamEvent(entities.UserCreatedEvent, 1, usr)
		if err != nil {
			logger.Error("error prepare message to broker", err)
			return "", err
		}
		if err = u.broker.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(entities.UserCreatedEvent),
			Value: b,
		}); err != nil {
			logger.Error("error stream event", err)
			return "", err
		}
		b, _ = json.Marshal(entities.UserBIEvent{
			PublicID: usr.PublicID,
		})
		if err = u.brokerBI.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(entities.UserBICreatedEvent),
			Value: b,
		}); err != nil {
			logger.Error("error stream event", err)
			return "", err
		}
	}
	//user exist or created, generate jwt
	return u.generateJWT(usr)
}

func (u *UserService) registerUser(googleUser *entities.GoogleUser) error {
	if err := u.uRepo.AddUser(googleUser, Worker); err != nil {
		logger.Error("error add user", err)
		return err
	}
	return nil
}

func (u *UserService) generateJWT(usr *user.UserAccountV1) (string, error) {
	atClaims := entities.UserJWT{
		usr.PublicID,
		usr.Version,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(100 * time.Hour).Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	jwtKey, err := at.SignedString(u.jwtKey)
	if err != nil {
		logger.Error("error generate jwt", err)
		return "", err
	}
	return jwtKey, nil
}

func (u *UserService) ChangeRole(publicID string, userVersion int64, newRole string) (*user.UserAccountV1, string, error) {
	usr, err := u.uRepo.ChangeRole(publicID, userVersion, newRole)
	if err != nil {
		logger.Error("error update user", err)
		return nil, "", err
	}

	// stream change event to broker
	b, _ := json.Marshal(usr)
	if err = u.broker.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(entities.UserUpdatedEvent),
		Value: b,
	}); err != nil {
		logger.Error("error stream event", err)
		return nil, "", err
	}

	jwtKey, err := u.generateJWT(usr)
	return usr, jwtKey, err
}

func (u *UserService) GetUserInfo(publicID string, userVersion int64) (*user.UserAccountV1, error) {
	return u.uRepo.GetByPublicID(publicID, userVersion)
}
