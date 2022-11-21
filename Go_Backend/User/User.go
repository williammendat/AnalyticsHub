package user

type (
	User struct {
		Firstname         string   `json:"Firstname" bson:"firstname"`
		Lastname          string   `json:"Lastname" bson:"lastname"`
		Email             string   `json:"Email" bson:"email"`
		Language          string   `json:"Language" bson:"language"`
		Base64ProfilePic  string   `json:"Base64ProfilePic" bson:"base65profilepic"`
		Favorites         []string `json:"Favorites" bson:"favorites"`
	}
)

func (user *User) CheckRequieredFields() {
	if user.Firstname == "" {
		panic("Firstname has to be provided")
	}

	if user.Lastname == "" {
		panic("Lastname has to be provided")
	}

	if user.Email == "" {
		panic("Email has to be provided")
	}

	if user.Language == "" {
		panic("Language has to be provided")
	}
}