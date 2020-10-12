package services

import (
	"github.com/abnergarcia1/GBM_test/pkg/gbm/data"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)


var investmentService = InvestmentService{}

func setEnv() {
	os.Setenv("OpenMarketHour", "1")
	os.Setenv("ClosedMarketHour", "23")
}

func TestCreateAccount(t *testing.T){
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}

	_,err:=investmentService.CreateAccount(newAccount)
	assert.Nil(t,err)

}

func TestInvestmentService_BuyOrder(t *testing.T) {
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 1,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	_, err:=investmentService.BuySellOrder(newOrder)
	assert.Nil(t, err)
}

func TestInvestmentService_BuyOrder_Balance(t *testing.T){
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 6,
		SharePrice: 5000,
	}

	newOrder.AccountID=resAccount.ID
	_, err:=investmentService.BuySellOrder(newOrder)
	if assert.NotNil(t, err){
		assert.Equal(t, "Insufficient Balance", err.Error(), "the balance is insufficient")
	}

}

func TestInvestmentService_SellOrder(t *testing.T){
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 6,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	_, _=investmentService.BuySellOrder(newOrder)

	newOrder.Operation="SELL"
	_, err:=investmentService.BuySellOrder(newOrder)
	assert.Nil(t, err)

}

func TestInvestmentService_SellOrder_Stocks(t *testing.T) {
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 6,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	_, _=investmentService.BuySellOrder(newOrder)

	newOrder.Operation="SELL"
	newOrder.TotalShares=60
	_, err:=investmentService.BuySellOrder(newOrder)
	if assert.NotNil(t, err){
		assert.Equal(t, "Insufficient Stocks", err.Error(), "The Ammount of stocks is insufficient")
	}
}

func TestInvestmentService_SellOrder_NotFoundIssuer(t *testing.T) {
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 6,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	_, _=investmentService.BuySellOrder(newOrder)

	newOrder.Operation="SELL"
	newOrder.IssuerName="NA"
	_, err:=investmentService.BuySellOrder(newOrder)
	if assert.NotNil(t, err){
		assert.Equal(t, "No Issuer Name was found in Account", err.Error(), "Issuer doesn't exist in account")
	}
}


func TestInvestmentService_Order_DuplicatedOperation(t *testing.T) {
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 1,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	_, _=investmentService.BuySellOrder(newOrder)
	time.Sleep(5 * time.Second)
	_, err:=investmentService.BuySellOrder(newOrder)
	if assert.NotNil(t, err){
		assert.Equal(t, "Duplicated Operation", err.Error(), "duplicated operation not permited")
	}
}

func TestInvestmentService_Order_AccountNotExists(t *testing.T) {
	setEnv()
	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	_, _= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 1,
		SharePrice: 50,
	}

	newOrder.AccountID=33
	_, err:=investmentService.BuySellOrder(newOrder)
	if assert.NotNil(t, err){
		assert.Equal(t, "Account not found", err.Error(), "Account doesn't exist")
	}
}

func TestInvestmentService_Order_ClosedMarket(t *testing.T){
	os.Setenv("OpenMarketHour", "2")
	os.Setenv("ClosedMarketHour", "3")

	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 1,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	_, err:=investmentService.BuySellOrder(newOrder)
	if assert.NotNil(t, err){
		assert.Equal(t, "Closed Market", err.Error(), "Market is closed, you cannot make any operation")
	}
}

func TestInvestmentService_BuySingleOrder_CorrectResponse(t *testing.T){
	setEnv()

	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 3,
		SharePrice: 50,
	}

	newOrder.AccountID=resAccount.ID
	response, _:=investmentService.BuySellOrder(newOrder)

	assert.Equal(t, response.CurrentBalance.Cash,newAccount.Cash - (newOrder.TotalShares*newOrder.SharePrice), "The account cash is wrong")
	assert.Equal(t, response.CurrentBalance.Issuers[0].TotalShares, newOrder.TotalShares, "Total Shares amount is wrong")
	assert.Equal(t, response.CurrentBalance.Issuers[0].IssuerName, newOrder.IssuerName, "Issuer name don't match")
	assert.Equal(t, response.CurrentBalance.Issuers[0].SharePrice, newOrder.SharePrice, "Share price don't match")
}

func TestInvestmentService_SellSingleOrder_CorrectResponse(t *testing.T){
	setEnv()

	investmentService.DB = &data.MemDB{}
	newAccount:=models.Account{
		Cash: 1000,
	}
	resAccount, _:= investmentService.CreateAccount(newAccount)

	newOrder:=models.Order{
		TimeStamp: time.Now().Unix(),
		Operation: "BUY",
		IssuerName: "AAPL",
		TotalShares: 3,
		SharePrice: 50,
	}
	newOrder.AccountID=resAccount.ID
	_, _=investmentService.BuySellOrder(newOrder)

	newOrder.Operation="SELL"
	newOrder.TotalShares=2

	response, _:=investmentService.BuySellOrder(newOrder)

	assert.Equal(t, response.CurrentBalance.Cash,newAccount.Cash - 50, "The account cash is wrong")
	assert.Equal(t, response.CurrentBalance.Issuers[0].TotalShares, int64(1), "Total Shares amount is wrong")
	assert.Equal(t, response.CurrentBalance.Issuers[0].IssuerName, newOrder.IssuerName, "Issuer name don't match")
	assert.Equal(t, response.CurrentBalance.Issuers[0].SharePrice, newOrder.SharePrice, "Share price don't match")
}