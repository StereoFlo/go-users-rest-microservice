package entity

import (
	"time"
)

type Token struct {
	ID                 int       `gorm:"primary_key;auto_increment" json:"id"`
	UserId             int       `gorm:"nut null" json:"user_id"`
	AccessToken        string    `gorm:"nut null;unique" json:"access_token"`
	RefreshToken       string    `gorm:"nut null;unique" json:"refresh_token"`
	AccessTokenExpire  time.Time `gorm:"nut null" json:"access_token_expire"`
	RefreshTokenExpire time.Time `gorm:"nut null" json:"refresh_token_expire"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}
