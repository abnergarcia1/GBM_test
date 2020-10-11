package services

import (
	"github.com/abnergarcia1/GBM_test/pkg/gbm/data"
	"github.com/abnergarcia1/GBM_test/pkg/gbm/models"
)

type InvestmentService struct{
	db data.IStorageData
}


func (s *InvestmentService) CreateAccount(account models.Account)(err error){

	s.db.CreateAccount()


}