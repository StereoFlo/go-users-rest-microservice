package persistence

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type UserRepo struct {
	db *gorm.DB
}

func CreateUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

var _ repository.UserRepository = &UserRepo{}

func (repo *UserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	dbErr := map[string]string{}
	err := repo.db.Debug().Update(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (repo *UserRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := repo.db.Debug().Preload("AccessTokens").Where("id = ?", id).Take(&user).Error
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
	err := repo.db.Debug().Preload("AccessTokens").Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *UserRepo) GetUserByAccessToken(token string) (*entity.Token, error) {
	var tokenEntity entity.Token
	err := repo.db.Debug().Where("access_token = ?", token).Take(&tokenEntity).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &tokenEntity, nil
}

func (repo *UserRepo) GetList() ([]entity.User, error) {
	var users []entity.User
	err := repo.db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}
