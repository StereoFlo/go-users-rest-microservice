package repository

import (
	"user-app/entity"
)

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetList() ([]entity.User, error)
	GetUserByAccessToken(token string) (*entity.User, error)
}
