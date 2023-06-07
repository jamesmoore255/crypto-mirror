package chainsync

import (
	"github.com/jamesmoore255/crypto-mirror/server/pkg/client"
	"github.com/jamesmoore255/crypto-mirror/server/pkg/logger"
)

// TODO: Create proper interface for this, could be cache, database, chain, etc
type WalletHoldingsFetcher interface {
	FetchWalletHoldings(bClient client.BlockClient)
}

type UserWallets struct {
	Addresses []string
}

// TODO: Do the prefix of the type here
func (uw UserWallets) FetchWalletHoldings(bClient client.BlockClient) {
	resultTokens := make(chan fetchWalletTokensResult)
	resultEth := make(chan fetchWalletEthResult)

	go fetchWalletTokens(uw.Addresses, bClient, resultTokens)
	go fetchWalletEth(uw.Addresses, bClient, resultEth)

	tokens := <-resultTokens
	eth := <-resultEth

	if tokens.Error != nil {
		logger.Errorf("Error fetching wallet tokens: %v", tokens.Error)
		panic(tokens.Error)
	}

	if eth.Error != nil {
		logger.Errorf("Error fetching wallet eth: %v", eth.Error)
		panic(eth.Error)
	}

	// Contract addresses for fetching token metadata
	var contractAddresses []string

	for _, tokenBalance := range tokens.Result.Result.TokenBalances {
		logger.Infof("Token: %v", tokenBalance)
		contractAddresses = append(contractAddresses, tokenBalance.ContractAddress)
	}

	logger.Infof("Contract addresses: %v", contractAddresses)

	resp, err := bClient.FetchMetadata(contractAddresses)

	if err != nil {
		logger.Errorf("Error fetching metadata: %v", err)
		return
	}

	logger.Infof("Metadata: %v", resp)

	// var walletTokens m.Wallet

	// config := &mapstructure.DecoderConfig{
	// 	DecodeHook: mapstructure.ComposeDecodeHookFunc(
	// 		utils.HexToDecimalHookFunc(),
	// 	),
	// 	Result: &walletTokens,
	// }

	// decoder, err := mapstructure.NewDecoder(config)
	// if err != nil {
	// 	logger.Errorf("Error creating decoder:", err)
	// 	return
	// }

	// err = decoder.Decode(tokens.Result.Result)

	// if err != nil {
	// 	logger.Errorf("Error decoding wallet tokens: %v", err)
	// 	return
	// }

	// logger.Infof("walletTokens:", walletTokens)
}

type fetchWalletTokensResult struct {
	Result client.FetchTokensResponse
	Error  error
}

func fetchWalletTokens(addresses []string, bClient client.BlockClient, resultChan chan fetchWalletTokensResult) {
	result, err := bClient.FetchTokens(addresses)

	if err != nil {
		logger.Errorf("Error fetching wallet tokens: %v", err)
		resultChan <- fetchWalletTokensResult{Error: err}
	}

	resultChan <- fetchWalletTokensResult{Result: result}
}

type fetchWalletEthResult struct {
	Result client.FetchEthResponse
	Error  error
}

func fetchWalletEth(addresses []string, bClient client.BlockClient, resultChan chan fetchWalletEthResult) {
	result, err := bClient.FetchEth(addresses)

	if err != nil {
		logger.Errorf("Error fetching wallet tokens: %v", err)
		resultChan <- fetchWalletEthResult{Error: err}
	}

	resultChan <- fetchWalletEthResult{Result: result}
}
