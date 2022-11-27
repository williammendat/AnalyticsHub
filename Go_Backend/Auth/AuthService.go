package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	user "github.com/williammendat/analytics-hub/User"
	utils "github.com/williammendat/analytics-hub/Utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service struct {
		tokenSecret string 
		userService *user.Service
	}
)

func NewService(userService *user.Service) *Service {

	tokenSecret := os.Getenv(utils.TokenSecret)
	if tokenSecret == "" {
		msg := fmt.Sprintf("%s is needed as env variable", tokenSecret)
		panic(msg)
	}

	return &Service{
		tokenSecret: tokenSecret,
		userService: userService,
	}
}

func (s *Service) SignUp(signUp SignUp) User {

	signUp.CheckRequieredFields()

	password := signUp.Password

	hashedPassword := hashAndSalt([]byte(password))

	newUser := user.User{
		UserBaseInfo: signUp.UserBaseInfo,
		Email: signUp.Email,
	}

	createdUser := s.userService.CreateUser(newUser, hashedPassword)
	
	token := s.createToken(createdUser.MongoID)

	result := User{
		User: createdUser.User,
		Token: token,
	}
	
	return result
}

func (s *Service) SignIn(signIn SignIn) User {

	signIn.CheckRequieredFields()

	foundUser := s.userService.GetUserByEmail(signIn.Email)

	validPassword := comparePasswords(foundUser.HashedPassword, []byte(signIn.Password))
	if !validPassword {
		panic("Password is wrong")
	}
	
	token := s.createToken(foundUser.MongoID)

	result := User{
		User: foundUser.User,
		Token: token,
	}
	
	return result
}

func (s *Service) createToken(userID primitive.ObjectID) string{
	ttl := 100 * time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		utils.ID: userID,
		"exp":  time.Now().UTC().Add(ttl).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		panic("Could not create Token string")
	}

	return tokenStr
}

func hashAndSalt(password []byte) string {
    
    // Use GenerateFromPassword to hash & salt pwd.
    // MinCost is just an integer constant provided by the bcrypt
    // package along with DefaultCost & MaxCost. 
    // The cost can be any value you want provided it isn't lower
    // than the MinCost (4)
    hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
    if err != nil {
        logrus.Println(err)
    }
    // GenerateFromPassword returns a byte slice so we need to
    // convert the bytes to a string and return it
    return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
    // Since we'll be getting the hashed password from the DB it
    // will be a string so we'll need to convert it to a byte slice
    byteHash := []byte(hashedPwd)
    err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
    if err != nil {
        logrus.Println(err)
        return false
    }
    
    return true
}