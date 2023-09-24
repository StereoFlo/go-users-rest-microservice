package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"user-app/internal/entity"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (repo *UserRepo) SaveUser(user *entity.User) (error, *entity.User) {
	err := repo.db.Debug().Create(&user).Error
	if err != nil {
		return err, nil
	}

	return nil, user
}

func (repo *UserRepo) GetUser(id int, tokenLimit int) (error, *entity.User) {
	var user entity.User
	err := repo.db.Debug().Preload("AccessTokens", func(db *gorm.DB) *gorm.DB {
		return db.Limit(tokenLimit)
	}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return err, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("user not found"), nil
	}

	return nil, &user
}

func (repo *UserRepo) GetUserByEmail(email string) (error, *entity.User) {
	var user entity.User
	err := repo.db.Debug().Where("email = ?", email).Take(&user).Error
	if err != nil {
		return err, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("user not found"), nil
	}

	return nil, &user
}

func (repo *UserRepo) GetList(limit int, offset int) (error, []entity.User) {
	var users []entity.User
	err := repo.db.Debug().Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return err, nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("user not found"), nil
	}

	return nil, users
}

func (repo *UserRepo) GetCount() (error, *int) {
	var cnt int
	err := repo.db.Table("users").Debug().Count(&cnt).Error
	if err != nil {
		return err, nil
	}

	return nil, &cnt
}
