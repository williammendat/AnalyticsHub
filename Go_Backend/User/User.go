package user

type (
	User struct {
		UserBaseInfo `json:",inline" bson:",inline"`
		Email             string   `json:"Email" bson:"email"`
		Base64ProfilePic  string   `json:"Base64ProfilePic" bson:"base65profilepic"`
		Favorites         []string `json:"Favorites" bson:"favorites"`
	}
)

func (user *User) CheckRequieredFields() {

	user.UserBaseInfo.CheckUserBaseInfos()

	if user.Email == "" {
		panic("Email has to be provided")
	}
}