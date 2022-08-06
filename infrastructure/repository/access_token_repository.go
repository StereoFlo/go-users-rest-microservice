package repository

import (
	"github.com/jinzhu/gorm"
	"user-app/entity"
)

type AccessTokenRepo struct {
	db *gorm.DB
}

func CreateTokenRepo(db *gorm.DB) *AccessTokenRepo {
	return &AccessTokenRepo{db}
}

func (repo AccessTokenRepo) SaveToken(token *entity.Token) (*entity.Token, map[string]string) {
	dbErr := map[string]string{}
	err := repo.db.Debug().Save(&token).Error
	if err != nil {
		return nil, dbErr
	}
	return token, nil
}
