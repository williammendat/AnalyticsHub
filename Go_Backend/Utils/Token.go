package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TokenSecret = "TokenSecret"
	ID = "id"
)

func ParseToken(jwtToken string) *jwt.Token {
	tokenSecret := os.Getenv(TokenSecret)
	if tokenSecret == "" {
		msg := fmt.Sprintf("%s is needed as env variable", tokenSecret)
		panic(msg)
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		panic("Could not parse Token")
	}

	return token
}

func GetUserID(jwtToken string) (primitive.ObjectID) {
	token := ParseToken(jwtToken)

	claims, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		panic("Could not get claims")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		panic("Could not get Id")
	}

	userMongoID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		panic("Could not parse string to object id")
	}
	
	return userMongoID
}