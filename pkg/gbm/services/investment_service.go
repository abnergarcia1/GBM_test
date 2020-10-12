package services

import (
	"errors"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/data"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"strconv"
	"time"
	"os"
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
	retAccount.Issuers=[]models.Order{}

	return

}

func (s *InvestmentService) BuySellOrder(order models.Order)(response models.OrderResponse, err error){
	response = models.OrderResponse{}
	query:=""

	account, err:=s.GetAccountDetails(order.AccountID)
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	response.BusinessErrors=[]string{}
	response.CurrentBalance=account

	response.CurrentBalance.Issuers=[]models.Order{}

	//check if the market is open
	err=s.IsOpenMarket()
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}


	//check if the operation is duplicated using 5 minutes tolerance

	err=operations.VerifyDuplicate(order)
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	//TODO: this can be replaced with transaction functions to ensure data integrity in db
	switch order.Operation{
	case "BUY":
		//check if the account has the enough balance
		err=s.HasEnoughBalance(account, order)
		if err!=nil{
			response.BusinessErrors=[]string{err.Error()}
			return
		}

		query="UPDATE accounts SET Cash=Cash - ? WHERE Id = ?"

		err=s.BuyShares(order)
		if err!=nil{
			response.BusinessErrors=[]string{err.Error()}
			return
		}

	case "SELL":
		//check if the account has enough stocks for issuer
		err=s.HasEnoughStocks(account, order)
		if err!=nil{
			response.BusinessErrors=[]string{err.Error()}
			return
		}

		query="UPDATE accounts SET Cash=Cash + ? WHERE Id = ?"

		err=s.SellShares(order)
		if err!=nil{
			response.BusinessErrors=[]string{err.Error()}
			return
		}

	default:
		err=errors.New("Invalid Operation Type")
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	if err = s.DB.Connect(); err != nil {
		response.BusinessErrors=[]string{err.Error()}
		return
	}
	defer s.DB.Disconnect()
	err=s.DB.Query(nil, query,order.TotalShares * order.SharePrice, account.ID)

	operations.AddOperation(order)

	account, err=s.GetAccountDetails(account.ID)
	if err!=nil{
		response.BusinessErrors=[]string{err.Error()}
		return
	}

	response.CurrentBalance=account
	response.CurrentBalance.ID=0

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

	return errors.New("No Issuer Name was found in Account")
}

func(s *InvestmentService) IsOpenMarket()(err error){
	openTime,_:=strconv.Atoi(os.Getenv("OpenMarketHour"))
	closeTime,_:=strconv.Atoi(os.Getenv("ClosedMarketHour"))
	hours, _,_:=time.Now().Clock()

	if hours<openTime || hours > closeTime{
		err=errors.New("Closed Market")
	}
	return
}

func(s *InvestmentService) GetAccountDetails(id int64)(account models.Account, err error){
	account=models.Account{}
	account.Issuers=[]models.Order{}
	var accountIssuers []models.Order

	if err = s.DB.Connect(); err != nil {
		return
	}
	defer s.DB.Disconnect()

	//Gets account main info
	err = s.DB.Query(&account, "SELECT Id, Cash FROM accounts WHERE Id=?", id)
	if err!=nil{
		return
	}

	if account.ID==0{
		err=errors.New("Account not found")
		return
	}

	//gets issuers to this account
	err = s.DB.Query(&accountIssuers, "SELECT IssuerName, TotalShares, SharePrice FROM orders WHERE AccountId=?", id)
	if err!=nil{
		return
	}

	account.Issuers=accountIssuers

	return
}

func(s *InvestmentService) BuyShares(order models.Order)(err error){

	if err = s.DB.Connect(); err != nil {
		return
	}
	defer s.DB.Disconnect()

	err=s.DB.Query(&order, "BUYSHARES", nil)

	return
}

func(s *InvestmentService) SellShares(order models.Order)(err error){
	if err = s.DB.Connect(); err != nil {
		return
	}
	defer s.DB.Disconnect()

	err=s.DB.Query(&order, "SELLSHARES", nil)

	return
}

