package models

type Order struct{
	AccountID	int64	`json:"-"`
	TimeStamp 	int64 	`json:"timestamp,omitempty"`
	Operation 	string  `json:"operation,omitempty"`
	IssuerName 	string 	`json:"issuer_name"`
	TotalShares int64 	`json:"total_shares"`
	SharePrice 	int64  	`json:"share_price"`
}

type OrderResponse struct{
	CurrentBalance Balance `json:"current_balance"`
	BusinessErrors []string `json:"business_errors"`

}