package auth

type (
	SignIn struct {
		Email string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
	}
)

func (s SignIn) CheckRequieredFields(){
	if s.Email == "" {
		panic("Email is requiered")
	}

	if s.Password == "" {
		panic("Password is requiered")
	}
}