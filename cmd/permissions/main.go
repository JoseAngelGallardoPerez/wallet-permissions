package main

import (
	"log"
	"runtime"

	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/app/config"
	"github.com/Confialink/wallet-permissions/internal/app/di"
	"github.com/Confialink/wallet-permissions/internal/routes"
)

var (
	appConfig *config.Main
)

func main() {
	c := di.Container
	appConfig = c.Config()

	if appConfig.Threads > 0 {
		runtime.GOMAXPROCS(appConfig.Threads)
		log.Printf("Running with %v threads", appConfig.Threads)
	}

	ginMode := env_mods.GetMode(appConfig.Env)
	gin.SetMode(ginMode)
	r := routes.GetRouter()

	log.Printf("Starting RPC on port: %s", appConfig.RPCPort)
	c.RPCServer().ListenAndServe(":" + appConfig.RPCPort)

	log.Printf("Starting API on port: %s", appConfig.Port)
	r.Run(":" + appConfig.Port)
}
