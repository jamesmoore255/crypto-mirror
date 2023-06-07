package client

import (
	"encoding/json"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/jamesmoore255/crypto-mirror/server/pkg/logger"
)

type Wallet struct {
	Address       string         `json:"address"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
	PageKey       string         `json:"pageKey"`
}

type TokenBalance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"` // hex
}

type BlockClient interface {
	FetchTokens(addresses []string) (FetchTokensResponse, error)
	FetchEth(addresses []string) (FetchEthResponse, error)
	FetchMetadata(contractAddresses []string) ([]FetchMetadataResponse, error)
}

type AlchemyAPI struct {
	URL string
}

type FetchTokensResponse struct {
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  Wallet `json:"result"`
}

func (ac AlchemyAPI) FetchTokens(addresses []string) (FetchTokensResponse, error) {
	client := resty.New()

	params := map[string]any{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenBalances",
		"params":  append(addresses, "erc20"),
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		Post(ac.URL)

	if err != nil {
		return FetchTokensResponse{}, err
	}

	// Handle the response
	var response FetchTokensResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return FetchTokensResponse{}, err
	}

	return response, nil
}

type FetchEthResponse struct {
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"` // hex
}

func (ac AlchemyAPI) FetchEth(addresses []string) (FetchEthResponse, error) {
	client := resty.New()

	params := map[string]any{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  append(addresses, "latest"),
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		Post(ac.URL)

	if err != nil {
		return FetchEthResponse{}, err
	}

	// Handle the response
	var response FetchEthResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return FetchEthResponse{}, err
	}

	return response, nil
}

type FetchMetadataResponse struct {
	ID      int64         `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Result  TokenMetadata `json:"result"`
}

type TokenMetadata struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logo"`
}

func (ac AlchemyAPI) FetchMetadata(contractAddresses []string) ([]FetchMetadataResponse, error) {
	var wg sync.WaitGroup
	wg.Add(len(contractAddresses))

	var mu sync.Mutex
	var tokenMetadata []FetchMetadataResponse

	for _, contractAddress := range contractAddresses {
		go func(addr string) {
			defer wg.Done()

			client := resty.New()

			params := map[string]any{
				"id":      1,
				"jsonrpc": "2.0",
				"method":  "alchemy_getTokenMetadata",
				"params":  []string{addr},
			}

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(params).
				Post(ac.URL)

			if err != nil {
				// Add to list of failed requests
				return
			}

			// Handle the response
			var response FetchMetadataResponse
			err = json.Unmarshal(resp.Body(), &response)
			if err != nil {
				// Add to list of failed requests
				logger.Errorf("Error unmarshalling response: %v", err)
				return
			}

			logger.Infof("Response: %v", response)
			mu.Lock()
			tokenMetadata = append(tokenMetadata, response)
			mu.Unlock()
		}(contractAddress)
	}

	wg.Wait()

	logger.Infof("Token metadata: %v", tokenMetadata)

	return tokenMetadata, nil
}
