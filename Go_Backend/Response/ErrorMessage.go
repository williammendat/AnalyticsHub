package response

type (
	ErrorMessage struct{
		Message string `json:"message" bson:"message"`
	}
)