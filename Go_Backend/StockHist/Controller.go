package stockhist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/williammendat/analytics-hub/Utils"
)

type (
	Controller struct {
		service *Service
	}
)

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) GetStockHist(g *gin.Context) {

	defer utils.BadRequestRecover(g)

	symbol := g.DefaultQuery("symbol", "")

	if symbol == "" {
		panic("symbol can not be null")
	}

	period := g.DefaultQuery("period", "1d")

	histFunc, ok := c.service.histFuncMap[period]
	if !ok {
		panic("Period is not supported")
	}

	histData := histFunc(symbol)

	g.IndentedJSON(http.StatusOK, histData)
}
