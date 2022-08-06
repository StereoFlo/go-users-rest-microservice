package application

import (
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type UserApp struct {
	UserRepo        repository.UserRepo
	AccessTokenRepo repository.AccessTokenRepo
}

func (userApp *UserApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return userApp.UserRepo.SaveUser(user)
}

func (userApp *UserApp) SaveToken(token *entity.Token) (*entity.Token, map[string]string) {
	return userApp.AccessTokenRepo.SaveToken(token)
}

func (userApp *UserApp) GetUser(userId int, tokenLimit int) (*entity.User, error) {
	return userApp.UserRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetUserByEmail(email string) (*entity.User, error) {
	return userApp.UserRepo.GetUserByEmail(email)
}

func (userApp *UserApp) GetUserByToken(userId int, tokenLimit int) (*entity.User, error) {
	return userApp.UserRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetList(limit int, offset int) ([]entity.User, error) {
	return userApp.UserRepo.GetList(limit, offset)
}

func (userApp *UserApp) GetUserByAccessToken(token string) (*entity.Token, error) {
	return userApp.UserRepo.GetUserByAccessToken(token)
}
