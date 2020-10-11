package models

type Account struct{
	ID 		int64 	`json:"id,omitempty"`
	Cash 	int64 	`json:"cash"`
	Issuers []Order `json:"issuers"`
}
