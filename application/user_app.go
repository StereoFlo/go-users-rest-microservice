package application

import (
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type UserApp struct {
	userRepo        *repository.UserRepo
	accessTokenRepo *repository.AccessTokenRepo
}

func NewUserApp(userRepo *repository.UserRepo, accessTokenRepo *repository.AccessTokenRepo) *UserApp {
	return &UserApp{userRepo: userRepo, accessTokenRepo: accessTokenRepo}
}

func (userApp *UserApp) SaveUser(user *entity.User) (*entity.User, error) {
	return userApp.userRepo.SaveUser(user)
}

func (userApp *UserApp) SaveToken(token *entity.Token) (*entity.Token, error) {
	return userApp.accessTokenRepo.SaveToken(token)
}

func (userApp *UserApp) GetTokenByUId(id string) (*entity.Token, error) {
	return userApp.accessTokenRepo.GetTokenByUId(id)
}

func (userApp *UserApp) GetUser(userId int, tokenLimit int) (*entity.User, error) {
	return userApp.userRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetUserByEmail(email string) (*entity.User, error) {
	return userApp.userRepo.GetUserByEmail(email)
}

func (userApp *UserApp) GetUserByToken(userId int, tokenLimit int) (*entity.User, error) {
	return userApp.userRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetList(limit int, offset int) ([]entity.User, error) {
	return userApp.userRepo.GetList(limit, offset)
}

func (userApp *UserApp) GetUserCount() (int, error) {
	return userApp.userRepo.GetCount()
}
