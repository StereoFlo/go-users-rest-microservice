package application

import (
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type UserApp struct {
	UserRepo        repository.UserRepository
	AccessTokenRepo repository.AccessTokenRepository
}

func (userApp *UserApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return userApp.SaveUser(user)
}

func (userApp *UserApp) SaveToken(token *entity.Token) (*entity.Token, map[string]string) {
	return userApp.AccessTokenRepo.SaveToken(token)
}

func (userApp *UserApp) GetUser(userId uint64) (*entity.User, error) {
	return userApp.UserRepo.GetUser(userId)
}

func (userApp *UserApp) GetUserByEmail(email string) (*entity.User, error) {
	return userApp.UserRepo.GetUserByEmail(email)
}

func (userApp *UserApp) GetUserByToken(userId uint64) (*entity.User, error) {
	return userApp.UserRepo.GetUser(userId)
}

func (userApp *UserApp) GetList() ([]entity.User, error) {
	return userApp.UserRepo.GetList()
}

func (userApp *UserApp) GetUserByAccessToken(token string) (*entity.Token, error) {
	return userApp.UserRepo.GetUserByAccessToken(token)
}
