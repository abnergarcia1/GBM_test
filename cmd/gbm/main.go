package main

import (
	"github.com/abnergarcia1/GBM_test/pkg/gbm"
)

func main(){

	api:=&gbm.APIsvc{}

	api.RunServer()
}