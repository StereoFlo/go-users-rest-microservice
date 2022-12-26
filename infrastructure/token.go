package infrastructure

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"os"
	"time"
)

type TokenData struct {
	UserId  int    `json:"user_id"`
	TokenId string `json:"token_id"`
}

type Claim struct {
	Data TokenData `json:"data"`
	jwt.RegisteredClaims
}

type Token struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewToken() (*Token, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_PRIVATE_KEY")))
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("JWT_PUBLIC_KEY")))
	if err != nil {
		return nil, err
	}

	return &Token{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (t Token) Get(ttl time.Time, userId int) (string, error) {
	now := time.Now()
	claims := make(jwt.MapClaims)
	uid := uuid.New()
	claims["data"] = TokenData{
		UserId:  userId,
		TokenId: uid.String(),
	}
	claims["exp"] = ttl.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (t Token) Validate(token string) (*Claim, error) {
	var c Claim
	_, err := jwt.ParseWithClaims(token, &c, t.parseToken(t.publicKey))
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (t Token) parseToken(key *rsa.PublicKey) func(jwtToken *jwt.Token) (interface{}, error) {
	return func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected method")
		}

		return key, nil
	}
}
