package data

import (
	"fmt"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"math/rand"
	"sync"
)

type MemDB struct {
	accountsTable []models.Account
	ordersTable []models.Order
	mux sync.Mutex
}


func(m *MemDB) Query (model interface{}, query string, args ...interface{}) (err error){

	switch query{
	case "INSERT INTO accounts(Cash) VALUES (?)":

		account:=models.Account{}
		account.Cash=args[0].(int64)
		account.ID=int64(rand.Intn(100000))

		m.mux.Lock()
		m.accountsTable = append(m.accountsTable, account)
		m.mux.Unlock()

		parseModel := model.(*models.Account)

		*parseModel=account

	case "SELECT IssuerName, TotalShares, SharePrice FROM orders WHERE AccountId=?":
		parseOrders := model.(*[]models.Order)
		*parseOrders=[]models.Order{}

		m.mux.Lock()
		for _, elem:=range m.ordersTable{
			if elem.AccountID == args[0].(int64){
				*parseOrders = append(*parseOrders, elem)
			}
		}
		m.mux.Unlock()

	case "SELECT Id, Cash FROM accounts WHERE Id=?":
		parseModel := model.(*models.Account)
		acctID:=args[0].(int64)

		for _, acct:=range m.accountsTable{

			if acct.ID==acctID{
				*parseModel=acct
				return
			}
		}


	case "UPDATE accounts SET Cash=Cash - ? WHERE Id = ?":
		cash:=args[0].(int64)
		acctID:=args[1].(int64)

		for i, acct:=range m.accountsTable{
			if acct.ID==acctID{
				m.accountsTable[i].Cash=acct.Cash-cash
				fmt.Println("acct buy info: ",acct)
				return
			}
		}

	case "UPDATE accounts SET Cash=Cash + ? WHERE Id = ?":
		cash:=args[0].(int64)
		acctID:=args[1].(int64)

		for _, acct:=range m.accountsTable{
			if acct.ID==acctID{
				acct.Cash=acct.Cash+cash
				return
			}
		}
	}

	fmt.Println("accounts table: ",m.accountsTable)
	fmt.Println("orders table: ", m.ordersTable)
	return
}


func (m *MemDB) Connect() error {
	return nil

}

func (m *MemDB) Disconnect()  {


}