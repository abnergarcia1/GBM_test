package data

import (
	"fmt"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
	"sync"
	"math/rand"
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

		for _, elem:=range m.ordersTable{
			if elem.AccountID == args[0].(int64){
				*parseOrders = append(*parseOrders, elem)
			}
		}


	}

	fmt.Println("the accounts table is :",m.accountsTable)

	return
}


func (m *MemDB) Connect() error {
	return nil

}

func (m *MemDB) Disconnect()  {


}