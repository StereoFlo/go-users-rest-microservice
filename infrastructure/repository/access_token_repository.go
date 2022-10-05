package repository

import (
	"github.com/jinzhu/gorm"
	"user-app/entity"
)

type AccessTokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *AccessTokenRepo {
	return &AccessTokenRepo{db}
}

func (repo AccessTokenRepo) SaveToken(token *entity.Token) (*entity.Token, error) {
	err := repo.db.Debug().Save(&token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (repo AccessTokenRepo) GetTokenByUId(id string) (*entity.Token, error) {
	var token entity.Token
	err := repo.db.Debug().Where("uuid = ?", id).Take(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}
