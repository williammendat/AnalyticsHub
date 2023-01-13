package stock

import "github.com/gin-gonic/gin"

type (
	Router struct {
		controller *Controller
	}
)

func NewRouter(service *Controller) *Router {
	return &Router{
		controller: service,
	}
}

func (r *Router) InitRouter(router *gin.Engine) {
	router.GET("/stock", r.controller.GetStockPage)
	router.GET("/stock/info", r.controller.GetStockData)
	router.GET("/stock/favs", r.controller.GetFavStockRanking)
}