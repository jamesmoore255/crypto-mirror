package balances

import (
	"errors"
	"strings"

	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/client"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/logger"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/model"
	utils "github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/utils"
	"github.com/shopspring/decimal"
)

type WalletsFetcher interface {
	FetchBalances(bClient client.BlockClient) (*[]WalletBalances, error)
}

type WalletAddresses []string

type ERC20WalletBalances *[]client.Wallet

type WalletBalances struct {
	Address string         `json:"address"`
	Tokens  []TokenBalance `json:"tokens"`
}

type TokenBalance struct {
	Token   model.TokenMetadata `json:"token"`
	Balance decimal.Decimal     `json:"balance"`
}

// http://localhost:8080/api/v1/wallets/balances?addresses=0x4587FC29677eCb64bCA501a7663D3C7ebEBCb27C
func (wa WalletAddresses) FetchBalances(bClient client.BlockClient) (*[]WalletBalances, error) {
	resultERC20Tokens := make(chan FetchWalletTokensResult)
	resultEtherBalances := make(chan FetchWalletEthResult)

	go wa.FetchWalletTokens(bClient, resultERC20Tokens)
	go wa.FetchWalletEth(bClient, resultEtherBalances)

	erc20Tokens := <-resultERC20Tokens
	etherBalances := <-resultEtherBalances

	if erc20Tokens.Error != nil {
		logger.Errorf("Error fetching wallet tokens: %v", erc20Tokens.Error)
		return &[]WalletBalances{}, errors.New("Error fetching tokens")
	}

	if etherBalances.Error != nil {
		logger.Errorf("Error fetching wallet eth: %v", etherBalances.Error)
		return &[]WalletBalances{}, errors.New("Error fetching eth")
	}

	metadataResp, err := FetchTokensMetadata(bClient, erc20Tokens.Result)

	if err != nil {
		logger.Errorf("Error fetching metadata: %v", err)
		return &[]WalletBalances{}, errors.New("Error fetching token metadata")
	}

	balances, err := parseBalances(erc20Tokens.Result, etherBalances.Result, metadataResp)

	if err != nil {
		return &[]WalletBalances{}, err
	}

	return &balances, nil
}

func parseBalances(erc20Tokens ERC20WalletBalances, ethBalances *[]client.EthBalance, metadataPtr *map[string]client.FetchMetadataResponse) ([]WalletBalances, error) {
	var balances []WalletBalances

	for _, wallet := range *erc20Tokens {
		walletRes := WalletBalances{Address: wallet.Address}

		for _, tokenBalance := range wallet.TokenBalances {

			metadataDeref := *metadataPtr

			metadata := metadataDeref[tokenBalance.ContractAddress].Result

			balance, err := utils.ConvertHexToDecimal(tokenBalance.TokenBalance, metadata.Decimals)

			if err != nil {
				logger.Errorf("Error converting hex to decimal: %v", err)
				return []WalletBalances{}, utils.GenericError()
			}

			walletRes.Tokens = append(walletRes.Tokens, TokenBalance{
				Balance: balance,
				Token:   metadata,
			})
		}

		balances = append(balances, walletRes)
	}

	for _, ethBalance := range *ethBalances {
		for i, wallet := range balances {
			if strings.ToLower(wallet.Address) == strings.ToLower(ethBalance.Address) {
				metadata := *utils.EthMetadata

				balance, err := utils.ConvertHexToDecimal(ethBalance.Balance, metadata.Decimals)

				if err != nil {
					logger.Errorf("Error converting hex to decimal: %v", err)
					return []WalletBalances{}, utils.GenericError()
				}

				balances[i].Tokens = append(balances[i].Tokens, TokenBalance{
					Balance: balance,
					Token:   metadata,
				})
				break
			}
		}
	}

	return balances, nil
}

func FetchTokensMetadata(bClient client.BlockClient, ws ERC20WalletBalances) (*map[string]client.FetchMetadataResponse, error) {
	contractAddresses := make(utils.StringSet)

	for _, wallet := range *ws {
		for _, tokenBalance := range wallet.TokenBalances {
			contractAddresses.Add(tokenBalance.ContractAddress)
		}
	}

	metadataResp, err := bClient.FetchMetadata(contractAddresses.Values())
	return metadataResp, err
}

type FetchWalletTokensResult struct {
	Result ERC20WalletBalances
	Error  error
}

func (wa WalletAddresses) FetchWalletTokens(bClient client.BlockClient, resultChan chan FetchWalletTokensResult) {
	result, err := bClient.FetchTokens(wa)

	if err != nil {
		logger.Errorf("Error fetching wallet tokens: %v", err)
		resultChan <- FetchWalletTokensResult{Error: err}
	}

	resultChan <- FetchWalletTokensResult{Result: ERC20WalletBalances(result)}
}

type FetchWalletEthResult struct {
	Result *[]client.EthBalance
	Error  error
}

func (wa WalletAddresses) FetchWalletEth(bClient client.BlockClient, resultChan chan FetchWalletEthResult) {
	result, err := bClient.FetchEth(wa)

	if err != nil {
		logger.Errorf("Error fetching wallet tokens: %v", err)
		resultChan <- FetchWalletEthResult{Error: err}
	}

	resultChan <- FetchWalletEthResult{Result: result}
}
