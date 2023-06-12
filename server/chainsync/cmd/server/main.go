package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	handler "github.com/jamesmoore255/crypto-mirror/server/chainsync/handlers"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/env"
	"github.com/jamesmoore255/crypto-mirror/server/chainsync/pkg/logger"
)

func main() {
	env.InitEnv()
	logger.InitLogger()

	r := gin.Default()

	handler.RegisterRoutes(r)

	// Set the address and port
	address := env.GetServerAddress()
	port := env.GetServerPort()
	addr := fmt.Sprintf("%s:%s", address, port)
	r.Run(addr) // listen and serve on

	logger.Info("Server started on at " + addr)
}
