package model

type TokenMetadata struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int32  `json:"decimals"`
	LogoURI  string `json:"logo"`
}
