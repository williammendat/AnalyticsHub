package stock

import (
	"net/http"
	"strconv"

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
func (c *Controller) GetStockPage(g *gin.Context) {
	
	defer utils.BadRequestRecover(g)
	limit, _ := strconv.Atoi(g.DefaultQuery("limit", "100"))
	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	term := g.DefaultQuery("term", "")

	pagedStocks := c.service.GetSearchPageStockModels(term, int64(limit), int64(page))

	g.IndentedJSON(http.StatusOK, pagedStocks)
}

func (c *Controller) GetStockData(g *gin.Context) {
	
	defer utils.BadRequestRecover(g)

	symbol := g.DefaultQuery("symbol", "")

	if symbol == "" {
		panic("symbol can not be null")
	}

	pagedStocks := c.service.GetStockData(symbol)

	g.IndentedJSON(http.StatusOK, pagedStocks)
}

func (c *Controller) GetFavStockRanking(g *gin.Context) {
	
	defer utils.BadRequestRecover(g)

	withHist, _ := strconv.ParseBool(g.DefaultQuery("withHist", "false"))

	favs := []string{"VRNS", "RPD", "AUDC", "DOMO"}

	stockFavs := c.service.GetStockFavorites(favs, withHist)

	g.IndentedJSON(http.StatusOK, stockFavs)
}


