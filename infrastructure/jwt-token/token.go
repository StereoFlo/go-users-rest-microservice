package jwt_token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

type Data struct {
	UserId  int    `json:"user_id"`
	TokenId string `json:"token_id"`
}

type Claim struct {
	Data Data `json:"data"`
	jwt.RegisteredClaims
}

type Token struct {
	privateKey []byte
	publicKey  []byte
}

func NewToken() Token {
	privateKey := getKeyData(os.Getenv("PRIVATE_KEY_FILE_PATH"))
	publicKey := getKeyData(os.Getenv("PUBLIC_KEY_FILE_PATH"))
	return Token{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (t Token) Get(ttl time.Time, userId int) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now()
	claims := make(jwt.MapClaims)
	uid := uuid.New()
	claims["data"] = Data{
		UserId:  userId,
		TokenId: uid.String(),
	}
	claims["exp"] = ttl.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (t Token) Validate(token string) (*Claim, error) {
	var c Claim
	key, err := jwt.ParseRSAPublicKeyFromPEM(t.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}
	_, err = jwt.ParseWithClaims(token, &c, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	return &c, nil
}

func getKeyData(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}
