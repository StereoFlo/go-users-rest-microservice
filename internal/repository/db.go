package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	entity2 "user-app/internal/entity"
)

type Repositories struct {
	User  *UserRepo
	Token *AccessTokenRepo
	db    *gorm.DB
}

func NewRepositories(DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //todo implement app_env
	})
	if err != nil {
		return nil, err
	}

	return &Repositories{
		User:  NewUserRepo(db),
		Token: NewTokenRepo(db),
		db:    db,
	}, nil
}

func (repo *Repositories) Automigrate() {
	err := repo.db.AutoMigrate(&entity2.User{}, &entity2.Token{})
	if err != nil {
		panic(err)
	}
}
