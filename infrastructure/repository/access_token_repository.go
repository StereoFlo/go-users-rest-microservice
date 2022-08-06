package repository

import "user-app/entity"

type AccessTokenRepository interface {
	SaveToken(token *entity.Token) (*entity.Token, map[string]string)
}
