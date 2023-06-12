package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/balances"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/client"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/env"
)

const v1PathPrefix = "/api/v1/"

func v1Path(suffix string) string {
	return fmt.Sprintf("%s%s", v1PathPrefix, suffix)
}

func RegisterRoutes(r *gin.Engine) {
	r.GET(v1Path("wallets/balances"), walletsTokensHandler)
}

//0x4587FC29677eCb64bCA501a7663D3C7ebEBCb27C

func walletsTokensHandler(c *gin.Context) {
	addresses := balances.WalletAddresses(c.QueryArray("addresses"))

	alchemyClient := client.AlchemyAPI{URL: env.GetAlchemyAPIURL()}

	holdings, err := addresses.FetchBalances(alchemyClient)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, holdings)
}
