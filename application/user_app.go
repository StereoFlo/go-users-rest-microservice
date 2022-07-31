package application

import (
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type userApp struct {
	userRepo repository.UserRepository
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetList() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByAccessToken(token string) (*entity.Token, error)
}

func (userApp *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return userApp.userRepo.SaveUser(user)
}

func (userApp *userApp) GetUser(userId uint64) (*entity.User, error) {
	return userApp.userRepo.GetUser(userId)
}

func (userApp *userApp) GetUserByEmail(email string) (*entity.User, error) {
	return userApp.userRepo.GetUserByEmail(email)
}

func (userApp *userApp) GetUserByToken(userId uint64) (*entity.User, error) {
	return userApp.userRepo.GetUser(userId)
}

func (userApp *userApp) GetList() ([]entity.User, error) {
	return userApp.userRepo.GetList()
}

func (userApp *userApp) GetUserByAccessToken(token string) (*entity.Token, error) {
	return userApp.userRepo.GetUserByAccessToken(token)
}
