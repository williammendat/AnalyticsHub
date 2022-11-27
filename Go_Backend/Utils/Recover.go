package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/williammendat/analytics-hub/Response"
)

func BadRequestRecover(g *gin.Context){
	if r := recover(); r != nil {
		err := fmt.Errorf("%v", r)

		errorMessage := response.ErrorMessage{
			Message: err.Error(),
		}

		g.IndentedJSON(http.StatusBadRequest, errorMessage)
	}
}