package entity

import (
	"time"
	"user-app/infrastructure/security"
)

type User struct {
	ID           uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Email        string    `gorm:"size:100;not null;unique" json:"email"`
	Password     string    `gorm:"size:100;not null;" json:"password"`
	AccessTokens []Token   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicUser struct {
	ID           uint64  `gorm:"primary_key;auto_increment" json:"id"`
	Email        string  `gorm:"size:100;not null;" json:"email"`
	AccessTokens []Token `json:"access_tokens"`
}

func (user *User) BeforeSave() error {
	hashPassword, err := security.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return nil
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

func (user *User) PublicUser() interface{} {
	return &PublicUser{
		ID:           user.ID,
		Email:        user.Email,
		AccessTokens: user.AccessTokens,
	}
}
