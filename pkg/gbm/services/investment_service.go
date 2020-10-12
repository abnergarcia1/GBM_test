package services

import (
	"errors"
	"fmt"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/data"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"time"
)

type InvestmentService struct{
	DB data.IStorageData
}

var operations  = models.StockOperationsQueue{}

func (s *InvestmentService) CreateAccount(account models.Account)(retAccount models.Account, err error){
	if err = s.DB.Connect(); err != nil {
		return
	}

	defer func() {
		s.DB.Disconnect()
	}()

	err=s.DB.Query(&retAccount, "INSERT INTO accounts(Cash) VALUES (?)",account.Cash)

	fmt.Println("investmentService val: ", retAccount)
	return

}

func (s *InvestmentService) BuySellOrder(order models.Order)(response models.OrderResponse, err error){
	response = models.OrderResponse{}

	//check if the market is open
	err=s.IsOpenMarket()
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	account:=models.Account{}
	var accountIssuers []models.Order

	if err = s.DB.Connect(); err != nil {
		response.BusinessErrors=[]string{err.Error()}
		return
	}
	defer s.DB.Disconnect()

	//check if the operation is duplicated using 5 minutes tolerance
	err=operations.VerifyDuplicate(order)

	//Gets account main info
	err = s.DB.Query(&account, "SELECT Id, Cash FROM accounts WHERE Id=?", order.AccountID)
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	//gets issuers to this account
	err = s.DB.Query(&accountIssuers, "SELECT IssuerName, TotalShares, SharePrice FROM orders WHERE AccountId=?", order.AccountID)
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	account.Issuers=accountIssuers

	query:=""
	switch order.Operation{
	case "BUY":
		//check if the account has the enough balance
		err=s.HasEnoughBalance(account, order)
		if err!=nil{
			response.BusinessErrors=[]string{err.Error()}
			return
		}

		query="UPDATE accounts SET Cash=Cash - ? WHERE Id = ?"

	case "SELL":
		//check if the account has enough stocks for issuer
		err=s.HasEnoughStocks(account, order)
		if err!=nil{
			response.BusinessErrors=[]string{err.Error()}
			return
		}

		query="UPDATE accounts SET Cash=Cash + ? WHERE Id = ?"

	}

	err=s.DB.Query(&account, query,account.Cash)

	//fmt.Println("investmentService val: ", retAccount)


	operations.AddOperation(order)

	return

}

func(s *InvestmentService) HasEnoughBalance(account models.Account, order models.Order)(err error){
	if account.Cash < order.TotalShares * order.SharePrice{
		err = errors.New("Insufficient Balance")
	}
	return
}

func(s *InvestmentService) HasEnoughStocks(account models.Account, order models.Order)(err error){
	for _, elem:=range account.Issuers{

		if elem.IssuerName==order.IssuerName {

			if elem.TotalShares<order.TotalShares{
				err=errors.New("Insufficient Stocks")
				return
			}else{
				return nil
			}

		}
	}

	return errors.New("No Issuer Name was found")
}

func(s *InvestmentService) IsOpenMarket()(err error){
	openTime:=6
	closeTime:=15
	hours, _,_:=time.Now().Clock()

	if hours<openTime || hours > closeTime{
		err=errors.New("Closed Market")
	}
	return
}