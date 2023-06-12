package client

import (
	"encoding/json"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/logger"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/model"
)

type Wallet struct {
	Address       string         `json:"address"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
	PageKey       string         `json:"pageKey"`
}

type EthBalance struct {
	Address string `json:"address"`
	Balance string `json:"balance"` // hex
}

type TokenBalance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"` // hex
}

type BlockClient interface {
	FetchTokens(addresses []string) ([]Wallet, error)
	FetchEth(addresses []string) ([]EthBalance, error)
	FetchMetadata(contractAddresses []string) (map[string]FetchMetadataResponse, error)
}

type AlchemyAPI struct {
	URL string
}

type TokenBalancesResponse struct {
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  Wallet `json:"result"`
}

func (ac AlchemyAPI) FetchTokens(addresses []string) ([]Wallet, error) {
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	var mu sync.Mutex

	client := resty.New()

	var result []Wallet

	for _, address := range addresses {
		go func(addr string) {
			defer wg.Done()

			params := map[string]any{
				"id":      1,
				"jsonrpc": "2.0",
				"method":  "alchemy_getTokenBalances",
				"params":  []string{addr, "erc20"},
			}

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(params).
				Post(ac.URL)

			if err != nil {
				return
			}

			// Handle the response
			var response TokenBalancesResponse
			err = json.Unmarshal(resp.Body(), &response)
			if err != nil {
				return
			}

			mu.Lock()
			result = append(result, response.Result)
			mu.Unlock()
		}(address)
	}

	wg.Wait()

	return result, nil
}

type FetchEthResponse struct {
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"` // hex
}

func (ac AlchemyAPI) FetchEth(addresses []string) ([]EthBalance, error) {
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	var mu sync.Mutex

	client := resty.New()

	var result []EthBalance

	for _, address := range addresses {
		go func(addr string) {
			defer wg.Done()

			params := map[string]any{
				"id":      1,
				"jsonrpc": "2.0",
				"method":  "eth_getBalance",
				"params":  []string{addr, "latest"},
			}

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(params).
				Post(ac.URL)

			if err != nil {
				return
			}

			// Handle the response
			var response FetchEthResponse
			err = json.Unmarshal(resp.Body(), &response)
			if err != nil {
				return
			}

			logger.Infof("Eth balance: %v", response)

			ethBalance := EthBalance{Address: addr, Balance: response.Result}

			mu.Lock()
			result = append(result, ethBalance)
			mu.Unlock()
		}(address)
	}

	wg.Wait()

	return result, nil
}

type FetchMetadataResponse struct {
	ID      int64               `json:"id"`
	JSONRPC string              `json:"jsonrpc"`
	Result  model.TokenMetadata `json:"result"`
}

func (ac AlchemyAPI) FetchMetadata(contractAddresses []string) (map[string]FetchMetadataResponse, error) {
	var wg sync.WaitGroup
	wg.Add(len(contractAddresses))

	var mu sync.Mutex
	tokenMetadata := make(map[string]FetchMetadataResponse)

	client := resty.New()

	for _, contractAddress := range contractAddresses {
		go func(addr string) {
			defer wg.Done()

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

			mu.Lock()
			tokenMetadata[addr] = response
			mu.Unlock()
		}(contractAddress)
	}

	wg.Wait()

	logger.Infof("Token metadata: %v", tokenMetadata)

	return tokenMetadata, nil
}
