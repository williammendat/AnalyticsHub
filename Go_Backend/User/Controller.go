package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/williammendat/analytics-hub/Utils"
)

type (
	Controller struct {
		userService *Service
	}
)

func NewController(userService *Service) *Controller {
	return &Controller{
		userService: userService,
	}
}

func (c *Controller) UpdateUser(g *gin.Context) {

	defer utils.BadRequestRecover(g)

	var user UserBaseInfo

	// Call BindJSON to bind the received JSON to
	// syncDates
	if err := g.BindJSON(&user); err != nil {

		panic(err)
	}

	//updatedUser := c.userService.UpdateUser()

	g.IndentedJSON(http.StatusOK, map[string]string{})
}