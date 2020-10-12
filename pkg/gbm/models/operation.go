

package models

import (
	"errors"
	"sync"

	"time"
)

type StockOperationsQueue struct{
	Operations []StockOperation
	mux sync.Mutex
}


type StockOperation struct{
	ID				int64
	Timestamp 		int64
	OperationType 	string
	Issuer 			string
	TotalShares 	int64
}

func(q *StockOperationsQueue) VerifyDuplicate(order Order) (err error){

		q.mux.Lock()
		defer q.mux.Unlock()

		for i, elem:=range q.Operations{

			if time.Now().Sub(time.Unix(elem.Timestamp,0)) < time.Minute*5 {
				//less than 5 min

				if elem.ID == order.AccountID &&
					elem.Issuer == order.IssuerName &&
					elem.OperationType == order.Operation &&
					elem.TotalShares == order.TotalShares {

					return errors.New("Duplicated Operation")
				}
			}else{
				//delete all operations great than 5 min
				q.Operations = append(q.Operations[:i],q.Operations[i+1:]...)

			}
		}

	return
}

func (q *StockOperationsQueue) AddOperation(order Order){
	q.mux.Lock()
	defer q.mux.Unlock()

	operation:=StockOperation{
		ID:order.AccountID,
		Timestamp: order.TimeStamp,
		OperationType: order.Operation,
		Issuer: order.IssuerName,
		TotalShares: order.TotalShares,
	}

	q.Operations = append(q.Operations, operation)

}
