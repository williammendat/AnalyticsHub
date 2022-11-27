package user

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const(
	DefaulBaseImage = "DefaulBaseImage"
)

type (
	Service struct {
		defaultImage string
		userRepository Repository
	}
)

func NewService(userRepository Repository) *Service {
	defaultImage := os.Getenv(DefaulBaseImage)
	if defaultImage == "" {
		panic("The default image variable is not set")
	}

	return &Service{
		defaultImage: defaultImage,
		userRepository: userRepository,
	}
}

func (service *Service) CreateUser(user User, hashedPassword string) UserModel {
	user.CheckRequieredFields()

	userToInsert := &UserModel{
		User: user,
		HashedPassword: hashedPassword,
	}

	userToInsert.Base64ProfilePic = service.defaultImage
	userToInsert.Favorites = make([]string, 0)

	newUser, err := service.userRepository.InsertOne(userToInsert)
	if err != nil {
		msg := fmt.Sprintf("Could not insert user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return newUser
}

func (service *Service) UpdateUser(Id primitive.ObjectID, oldUser UserBaseInfo) User {

	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	hasChanges := false

	if oldUser.Firstname != "" {
		user.Firstname = oldUser.Firstname
		hasChanges = true
	}

	if oldUser.Lastname != "" {
		user.Firstname = oldUser.Lastname
		hasChanges = true
	}

	if oldUser.Language != "" {
		user.Firstname = oldUser.Language
		hasChanges = true
	}

	if !hasChanges {
		return user.User
	}

	updatedUser, err := service.userRepository.ReplaceOne(filter, &user)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return updatedUser.User
}

func (service *Service) UpdateUserProfilePicture(Id primitive.ObjectID, profileImage string) User {

	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	user.Base64ProfilePic = profileImage

	updatedUser, err := service.userRepository.ReplaceOne(filter, &user)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return updatedUser.User
}

func (service *Service) AddUserFavourite(Id primitive.ObjectID, symbol string) User {

	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	if len(user.Favorites) == 0 {
		user.Favorites = make([]string, 0, 1)
	}

	if len(user.Favorites) >= 5 {
		panic("User can not have more than with Favourites")
	}

	user.Favorites = append(user.Favorites, symbol)

	updatedUser, err := service.userRepository.ReplaceOne(filter, &user)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return updatedUser.User
}

func (service *Service) RemoveUserFavourite(Id primitive.ObjectID, symbol string) User {

	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	if len(user.Favorites) == 0 {
		return user.User
	}

	newFavs := make([]string, 0, len(user.Favorites) - 1)
	for _, fav := range user.Favorites {
		if fav == symbol {
			continue
		}
		newFavs = append(newFavs, fav)
	}

	user.Favorites = newFavs

	updatedUser, err := service.userRepository.ReplaceOne(filter, &user)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return updatedUser.User
}

func (service *Service) GetUser(Id primitive.ObjectID) UserModel {
	filter := mongodb.NewFilter(mongodb.WithMongoID(Id))

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic(msg)
	}

	return user
}

func (service *Service) GetUserByEmail(email string) UserModel {
	filter := bson.M{"email": email}

	user, err := service.userRepository.FindOne(filter)
	if err != nil {
		msg := fmt.Sprintf("Could not get user: %s", err)
		logrus.Errorln(msg)
		panic("Could not find user with email: " + email)
	}

	return user
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
