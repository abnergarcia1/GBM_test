package gbm

import (
	"fmt"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/services"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var operations  = models.StockOperationsQueue{}

type APIHandlers struct{
	investmentService services.InvestmentService
}

func (h *APIHandlers) CreateAccount(c *gin.Context){


}

func (h *APIHandlers) BuySellOrder(c *gin.Context){
	id:=c.Param("id")

	message:="user id is " + id

	order:=models.Order{}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	order.AccountID,_=strconv.ParseInt(id,10,64)
	err:=operations.VerifyDuplicate(order)
	if err!=nil{
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(operations.Operations)

	c.String(http.StatusOK,message)

}
