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
				elem.Operation=""
				elem.TimeStamp=0
				*parseOrders = append(*parseOrders, elem)
			}
		}
		m.mux.Unlock()


	case "SELECT Id, Cash FROM accounts WHERE Id=?":
		parseModel := model.(*models.Account)
		acctID:=args[0].(int64)

		m.mux.Lock()
		for _, acct:=range m.accountsTable{
			if acct.ID==acctID{
				*parseModel=acct
				break
			}
		}
		m.mux.Unlock()

	case "UPDATE accounts SET Cash=Cash - ? WHERE Id = ?":
		cash:=args[0].(int64)
		acctID:=args[1].(int64)

		m.mux.Lock()
		for i, acct:=range m.accountsTable{
			if acct.ID==acctID{
				m.accountsTable[i].Cash=acct.Cash-cash
				fmt.Println("acct buy info: ",acct)
				break
			}
		}
		m.mux.Unlock()

	case "UPDATE accounts SET Cash=Cash + ? WHERE Id = ?":
		cash:=args[0].(int64)
		acctID:=args[1].(int64)

		m.mux.Lock()
		for i, acct:=range m.accountsTable{
			if acct.ID==acctID{
				m.accountsTable[i].Cash=acct.Cash+cash
				fmt.Println("acct buy info: ",acct)
				break
			}
		}
		m.mux.Unlock()

	case "BUYSHARES":
		parseBuyOrder := model.(*models.Order)

		for i, order := range m.ordersTable{

			if parseBuyOrder.AccountID==order.AccountID && parseBuyOrder.IssuerName==order.IssuerName{

				m.mux.Lock()
				m.ordersTable[i].TotalShares+=parseBuyOrder.TotalShares
				m.mux.Unlock()
				return
			}
		}

		m.ordersTable = append(m.ordersTable, *parseBuyOrder)

	case "SELLSHARES":
		parseSellOrder := model.(*models.Order)

		for i, order := range m.ordersTable{

			if parseSellOrder.AccountID==order.AccountID && parseSellOrder.IssuerName==order.IssuerName{

				m.mux.Lock()
				m.ordersTable[i].TotalShares-=parseSellOrder.TotalShares
				if m.ordersTable[i].TotalShares<=0{
					m.ordersTable = append(m.ordersTable[:i],m.ordersTable[i+1:]...)
				}
				m.mux.Unlock()
				return
			}
		}

	}

	return
}


func (m *MemDB) Connect() error {
	return nil

}

func (m *MemDB) Disconnect()  {


}