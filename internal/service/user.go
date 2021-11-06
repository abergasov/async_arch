package service

import (
	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/repository/user"
)

const (
	Worker    = "worker"
	AdminRole = "admin"
)

type UserService struct {
	uRepo user.UserRepo
}

func InitUserService(uRepo user.UserRepo) *UserService {
	return &UserService{uRepo: uRepo}
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
	}
	//user exist or created, generate jwt
	return "132312", nil
}

func (u *UserService) registerUser(googleUser *entities.GoogleUser) error {
	if err := u.uRepo.AddUser(googleUser, Worker); err != nil {
		logger.Error("error add user", err)
		return err
	}
	return nil
}
