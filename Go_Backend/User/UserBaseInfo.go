package user

type (
	UserBaseInfo struct {
		Firstname         string   `json:"Firstname" bson:"firstname"`
		Lastname          string   `json:"Lastname" bson:"lastname"`
		Language          string   `json:"Language" bson:"language"`
	}
)

func (u UserBaseInfo) CheckUserBaseInfos() {
	if u.Firstname == "" {
		panic("Firstname has to be provided")
	}

	if u.Lastname == "" {
		panic("Lastname has to be provided")
	}

	if u.Language == "" {
		panic("Language has to be provided")
	}
}