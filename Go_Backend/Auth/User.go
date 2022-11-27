package auth

import user "github.com/williammendat/analytics-hub/User"

type (
	User struct {
		user.User `json:",inline" bson:",inline"`
		Token string `json:"token" bson:"token"`
	}
)