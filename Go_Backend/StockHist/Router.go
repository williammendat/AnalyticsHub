package stockhist

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
	router.GET("/stockHist", r.controller.GetStockHist)
}