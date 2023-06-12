package model

import "github.com/shopspring/decimal"

type Wallet struct {
	Address       string         `json:"address"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
	PageKey       string         `json:"pageKey"`
}

type TokenBalance struct {
	ContractAddress string          `json:"contractAddress"`
	TokenBalance    decimal.Decimal `json:"tokenBalance"`
}
