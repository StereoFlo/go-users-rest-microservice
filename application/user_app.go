package application

import (
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type UserApp struct {
	UserRepo        *repository.UserRepo
	AccessTokenRepo *repository.AccessTokenRepo
}

func NewUserApp(userRepo *repository.UserRepo, accessTokenRepo *repository.AccessTokenRepo) *UserApp {
	return &UserApp{UserRepo: userRepo, AccessTokenRepo: accessTokenRepo}
}

func (userApp *UserApp) SaveUser(user *entity.User) (*entity.User, error) {
	return userApp.UserRepo.SaveUser(user)
}

func (userApp *UserApp) SaveToken(token *entity.Token) (*entity.Token, error) {
	return userApp.AccessTokenRepo.SaveToken(token)
}

func (userApp *UserApp) GetTokenByUId(id string) (*entity.Token, error) {
	return userApp.AccessTokenRepo.GetTokenByUId(id)
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

func (userApp *UserApp) GetUserCount() (int, error) {
	return userApp.UserRepo.GetCount()
}
