package repository

import (
	"github.com/jinzhu/gorm"
	"user-app/internal/entity"
)

type AccessTokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *AccessTokenRepo {
	return &AccessTokenRepo{db}
}

func (repo AccessTokenRepo) SaveToken(token *entity.Token) (error, *entity.Token) {
	err := repo.db.Debug().Save(&token).Error
	if err != nil {
		return err, nil
	}
	return nil, token
}

func (repo AccessTokenRepo) GetTokenByUId(id string) (error, *entity.Token) {
	var token entity.Token
	err := repo.db.Debug().Where("uuid = ?", id).Take(&token).Error
	if err != nil {
		return err, nil
	}
	return nil, &token
}
