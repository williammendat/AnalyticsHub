package auth

import user "github.com/williammendat/analytics-hub/User"

type (
	SignUp struct {
		SignIn `json:",inline" bson:",inline"`
		user.UserBaseInfo `json:",inline" bson:",inline"`
	}
)

func (s SignUp) CheckRequieredFields() {
	s.SignIn.CheckRequieredFields()
	s.CheckUserBaseInfos()
}