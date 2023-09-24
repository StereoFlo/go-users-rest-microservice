package application

import (
	entity2 "user-app/internal/entity"
	repository2 "user-app/internal/infrastructure/repository"
)

type UserApp struct {
	userRepo        *repository2.UserRepo
	accessTokenRepo *repository2.AccessTokenRepo
}

func NewUserApp(userRepo *repository2.UserRepo, accessTokenRepo *repository2.AccessTokenRepo) *UserApp {
	return &UserApp{userRepo: userRepo, accessTokenRepo: accessTokenRepo}
}

func (userApp *UserApp) SaveUser(user *entity2.User) (*entity2.User, error) {
	return userApp.userRepo.SaveUser(user)
}

func (userApp *UserApp) SaveToken(token *entity2.Token) (*entity2.Token, error) {
	return userApp.accessTokenRepo.SaveToken(token)
}

func (userApp *UserApp) GetTokenByUId(id string) (*entity2.Token, error) {
	return userApp.accessTokenRepo.GetTokenByUId(id)
}

func (userApp *UserApp) GetUser(userId int, tokenLimit int) (*entity2.User, error) {
	return userApp.userRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetUserByEmail(email string) (*entity2.User, error) {
	return userApp.userRepo.GetUserByEmail(email)
}

func (userApp *UserApp) GetUserByToken(userId int, tokenLimit int) (*entity2.User, error) {
	return userApp.userRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetList(limit int, offset int) ([]entity2.User, error) {
	return userApp.userRepo.GetList(limit, offset)
}

func (userApp *UserApp) GetUserCount() (int, error) {
	return userApp.userRepo.GetCount()
}
