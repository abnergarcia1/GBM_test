package gbm

import (
	"github.com/gin-gonic/gin"
	"log"
)

var (
	router = gin.Default()
)

type APIsvc struct{}

func (a *APIsvc) RunServer() {

	handlers := APIHandlers{}
	router.POST("/accounts", handlers.CreateAccount)
	router.POST("/accounts/:id/orders",handlers.BuySellOrder)
	log.Fatal(router.Run(":8085"))

}