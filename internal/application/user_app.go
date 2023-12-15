package application

import (
	entity2 "user-app/internal/entity"
	"user-app/internal/repository"
)

type UserApp struct {
	userRepo        *repository.UserRepo
	accessTokenRepo *repository.AccessTokenRepo
}

func NewUserApp(userRepo *repository.UserRepo, accessTokenRepo *repository.AccessTokenRepo) *UserApp {
	return &UserApp{userRepo: userRepo, accessTokenRepo: accessTokenRepo}
}

func (userApp *UserApp) SaveUser(user *entity2.User) (error, *entity2.User) {
	return userApp.userRepo.SaveUser(user)
}

func (userApp *UserApp) SaveToken(token *entity2.Token) (error, *entity2.Token) {
	return userApp.accessTokenRepo.SaveToken(token)
}

func (userApp *UserApp) GetTokenByUId(id string) (error, *entity2.Token) {
	return userApp.accessTokenRepo.GetTokenByUId(id)
}

func (userApp *UserApp) GetUser(userId int, tokenLimit int) (error, *entity2.User) {
	return userApp.userRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetUserByEmail(email string) (error, *entity2.User) {
	return userApp.userRepo.GetUserByEmail(email)
}

func (userApp *UserApp) GetUserByToken(userId int, tokenLimit int) (error, *entity2.User) {
	return userApp.userRepo.GetUser(userId, tokenLimit)
}

func (userApp *UserApp) GetList(limit int, offset int) (error, []entity2.User) {
	return userApp.userRepo.GetList(limit, offset)
}

func (userApp *UserApp) GetUserCount() (error, *int) {
	return userApp.userRepo.GetCount()
}
