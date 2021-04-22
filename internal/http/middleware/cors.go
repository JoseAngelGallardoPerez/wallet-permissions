package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-pkg-env_config"
)

func Cors(config *env_config.Cors) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowMethods = config.Methods
	for _, origin := range config.Origins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
		}
	}
	if !corsConfig.AllowAllOrigins {
		corsConfig.AllowOrigins = config.Origins
	}
	corsConfig.AllowHeaders = config.Headers

	return cors.New(corsConfig)
}
