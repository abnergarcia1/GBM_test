package gbm

import (
	"github.com/abnergarcia1/GBM_test/pkg/gbm/data"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	router = gin.Default()
)

type APIsvc struct{}

func (a *APIsvc) RunServer() {

	handlers := APIHandlers{}
	handlers.investmentService.DB = &data.MemDB{}

	router.POST("/accounts", handlers.CreateAccount)
	router.POST("/accounts/:id/orders",handlers.BuySellOrder)
	log.Fatal(router.Run(":8085"))

}