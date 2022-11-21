package user

import (
	"fmt"

	"github.com/sirupsen/logrus"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Service struct {
		userRepository Repository
	}
)

func NewService(userRepository Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (service *Service) CreateUser(user UserModel) User {
	user.CheckRequieredFields()

	user.Base64ProfilePic = "Lol"
	user.Favorites = make([]string, 0)
	user.Favorites = append(user.Favorites, "AAPL")

	newUser, err := service.userRepository.InsertOne(&user)
	if err != nil {
		msg := fmt.Sprintf("Could not insert user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return newUser.User
}

func (service *Service) UpdateUser(Id primitive.ObjectID, oldUser User) User {
	oldUser.CheckRequieredFields()

	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	user.Firstname = oldUser.Firstname
	user.Lastname = oldUser.Lastname
	user.Email = oldUser.Email
	user.Language = oldUser.Language
	user.Base64ProfilePic = oldUser.Base64ProfilePic
	user.Favorites = oldUser.Favorites

	updatedUser, err := service.userRepository.ReplaceOne(filter, &user)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return updatedUser.User
}

func (service *Service) GetUser(Id primitive.ObjectID) User {
	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return user.User
}

func (service *Service) DeleteUser(Id primitive.ObjectID) {
	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	err := service.userRepository.DeleteOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not delete user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}
}
