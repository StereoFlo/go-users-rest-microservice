package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"user-app/internal/entity"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (repo *UserRepo) SaveUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := repo.db.Debug().WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepo) GetUser(ctx context.Context, id int, tokenLimit int) (*entity.User, error) {
	var user entity.User
	err := repo.db.Debug().WithContext(ctx).Preload("Tokens", func(db *gorm.DB) *gorm.DB {
		return db.Limit(tokenLimit)
	}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := repo.db.Debug().WithContext(ctx).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (repo *UserRepo) GetList(ctx context.Context, limit int, offset int) ([]entity.User, error) {
	var users []entity.User
	err := repo.db.WithContext(ctx).Debug().Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return users, nil
}

func (repo *UserRepo) GetCount(cxt context.Context) (error, *int64) {
	var cnt int64
	err := repo.db.WithContext(cxt).Table("users").Debug().Count(&cnt).Error
	if err != nil {
		return err, nil
	}

	return nil, &cnt
}
