package models

type Balance struct{
	Cash 	int64	 `json:"cash,omitempty"`
	Issuers []Order  `json:"issuers,omitempty"`
}
