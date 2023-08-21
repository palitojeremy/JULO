package model

import (
	"time"
)

type Wallet struct {
	Id        		string    `json:"id" gorm:"primary_key"`
	OwnedBy   		string    `json:"owned_by"`
	Status    		string    `json:"status"`
	EnabledAt 		time.Time `json:"enabled_at"`
	Balance				int32			`json:"balance"`
}

type Transaction struct {
	Id        		string    `json:"id" gorm:"primary_key"`
	Type   				string    `json:"type"`
	Status    		string    `json:"status"`
	TransactedBy	string		`json:"transacted_by"`
	TransactedAt 	time.Time `json:"transacted_at"`
	Amount				int32			`json:"amount"`
	ReferenceId 	string		`json:"reference_id"`
}

type TransactionInput struct {
	Amount        int32     `json:"amount"`
	ReferenceId   string    `json:"reference_id"`
}

type CreateWalletInput struct {
	Id		        string    `json:"customer_xid"`
}

type DisableWalletInput struct {
	IsDisabled		bool    `json:"is_disabled"`
}