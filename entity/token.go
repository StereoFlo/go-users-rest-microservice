package entity

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Token struct {
	gorm.Model
	UserId             int
	AccessToken        string    `gorm:"size:36;nut null" json:"access_token"`
	RefreshToken       string    `gorm:"size:36;nut null" json:"refresh_token"`
	AccessTokenExpire  time.Time `gorm:"nut null" json:"access_token_expire"`
	RefreshTokenExpire time.Time `gorm:"nut null" json:"refresh_token_expire"`
}
