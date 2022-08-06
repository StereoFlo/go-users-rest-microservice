package persistence

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"user-app/entity"
	"user-app/infrastructure/repository"
)

type Repositories struct {
	User  repository.UserRepository
	Token repository.AccessTokenRepository
	db    *gorm.DB
}

func CreateRepositories(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	dBUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, dBUrl)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User:  CreateUserRepo(db),
		Token: CreateTokenRepo(db),
		db:    db,
	}, nil
}

func (repo *Repositories) Close() error {
	return repo.db.Close()
}

func (repo *Repositories) Automigrate() error {
	return repo.db.AutoMigrate(&entity.User{}, &entity.Token{}).Error
}
