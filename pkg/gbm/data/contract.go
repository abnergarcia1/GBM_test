package data

import(
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
)

type IStorageData interface {
	GetCurrentBalance(id string) (currentBalance models.Balance, err error)
	GetIssuersById(id string)(issuers []models.Order, err error)
	CreateAccount(account models.Account) (err error)
	SellOrder(order models.Order) (err error)
	BuyOrder(order models.Order) (err error)
}
