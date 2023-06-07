package main

import (
	"github.com/jamesmoore255/crypto-mirror/server/internal/chainsync"
	"github.com/jamesmoore255/crypto-mirror/server/pkg/client"
	"github.com/jamesmoore255/crypto-mirror/server/pkg/env"
	"github.com/jamesmoore255/crypto-mirror/server/pkg/logger"
)

func main() {
	env.InitEnv()
	logger.InitLogger()
	alchemyClient := client.AlchemyAPI{URL: env.GetAlchemyAPIURL()}
	chainsync.UserWallets{Addresses: []string{"0x4587FC29677eCb64bCA501a7663D3C7ebEBCb27C"}}.FetchWalletHoldings(alchemyClient)
}
