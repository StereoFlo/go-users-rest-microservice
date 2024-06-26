package application

import (
	"context"
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

func (userApp *UserApp) SaveUser(ctx context.Context, user *entity2.User) (*entity2.User, error) {
	return userApp.userRepo.SaveUser(ctx, user)
}

func (userApp *UserApp) SaveToken(ctx context.Context, token *entity2.Token) (*entity2.Token, error) {
	return userApp.accessTokenRepo.SaveToken(ctx, token)
}

func (userApp *UserApp) GetTokenByUId(id string) (*entity2.Token, error) {
	return userApp.accessTokenRepo.GetTokenByUId(id)
}

func (userApp *UserApp) GetUser(ctx context.Context, userId int, tokenLimit int) (*entity2.User, error) {
	return userApp.userRepo.GetUser(ctx, userId, tokenLimit)
}

func (userApp *UserApp) GetUserByEmail(ctx context.Context, email string) (*entity2.User, error) {
	return userApp.userRepo.GetUserByEmail(ctx, email)
}

func (userApp *UserApp) GetUserByToken(ctx context.Context, userId int, tokenLimit int) (*entity2.User, error) {
	return userApp.userRepo.GetUser(ctx, userId, tokenLimit)
}

func (userApp *UserApp) GetList(ctx context.Context, limit int, offset int) ([]entity2.User, error) {
	return userApp.userRepo.GetList(ctx, limit, offset)
}

func (userApp *UserApp) GetUserCount(cxt context.Context) (error, *int64) {
	return userApp.userRepo.GetCount(cxt)
}
