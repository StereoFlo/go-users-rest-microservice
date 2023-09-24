package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	entity2 "user-app/internal/entity"
)

type Repositories struct {
	User  *UserRepo
	Token *AccessTokenRepo
	db    *gorm.DB
}

func NewRepositories(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	dBUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, dBUrl)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User:  NewUserRepo(db),
		Token: NewTokenRepo(db),
		db:    db,
	}, nil
}

func (repo *Repositories) Close() error {
	return repo.db.Close()
}

func (repo *Repositories) Automigrate() error {
	return repo.db.AutoMigrate(&entity2.User{}, &entity2.Token{}).Error
}
