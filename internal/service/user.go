package service

import (
	"time"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/repository/user"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

const (
	Worker    = "worker"
	AdminRole = "admin"
)

type UserService struct {
	uRepo  user.UserRepo
	jwtKey []byte
}

func InitUserService(uRepo user.UserRepo, jwtKey string) *UserService {
	return &UserService{uRepo: uRepo, jwtKey: []byte(jwtKey)}
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

func (u *UserService) generateJWT(usr *entities.UserAccount) (string, error) {
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

func (u *UserService) ChangeRole(publicID uuid.UUID, userVersion int, newRole string) (*entities.UserAccount, string, error) {
	usr, err := u.uRepo.ChangeRole(publicID, userVersion, newRole)
	if err != nil {
		logger.Error("error update user", err)
		return nil, "", err
	}
	jwtKey, err := u.generateJWT(usr)
	return usr, jwtKey, err
}

func (u *UserService) GetUserInfo(publicID uuid.UUID, userVersion int) (*entities.UserAccount, error) {
	return u.uRepo.GetUserByPublicID(publicID, userVersion)
}
