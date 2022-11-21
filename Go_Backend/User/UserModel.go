package user

import mongodb "github.com/williammendat/analytics-hub/MongoDb"

type (
	UserModel struct {
		mongodb.BaseModel `json:",inline" bson:",inline"`
		User              `json:",inline" bson:",inline"`
		HashedPassword string `json:"HashedPassword" bson:"hashedpassword"`
	}
)

func (u UserModel) GetDoc() UserModel {
	return u
}
