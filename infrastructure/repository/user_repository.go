package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
	"user-app/entity"
)

type UserRepo struct {
	Database *gorm.DB
}

func CreateUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (repo *UserRepo) SaveUser(user *entity.User) (*entity.User, error) {
	err := repo.Database.Debug().Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

func (repo *UserRepo) GetUser(id int, tokenLimit int) (*entity.User, error) {
	var user entity.User
	err := repo.Database.Debug().Preload("AccessTokens", func(db *gorm.DB) *gorm.DB {
		return db.Limit(tokenLimit)
	}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := repo.Database.Debug().Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *UserRepo) GetList(limit int, offset int) ([]entity.User, error) {
	var users []entity.User
	err := repo.Database.Debug().Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (repo *UserRepo) GetCount() (int, error) {
	var cnt int
	err := repo.Database.Table("users").Debug().Count(&cnt).Error
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}
