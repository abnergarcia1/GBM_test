package gbm

import (
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)



type APIHandlers struct{
	investmentService services.InvestmentService
}

func (h *APIHandlers) CreateAccount(c *gin.Context){

	account:=models.Account{}

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	retacct, err:=h.investmentService.CreateAccount(account)
	if err!=nil{
		c.String(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, retacct)

}

func (h *APIHandlers) BuySellOrder(c *gin.Context){
	id:=c.Param("id")

	order:=models.Order{}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	order.AccountID,_=strconv.ParseInt(id,10,64)

	orderResponse, err:=h.investmentService.BuySellOrder(order)
	if err!=nil{
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK,orderResponse)

}
