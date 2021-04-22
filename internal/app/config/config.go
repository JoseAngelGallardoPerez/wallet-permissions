package config

import (
	"github.com/Confialink/wallet-pkg-env_config"
)

type Main struct {
	Env     string
	Threads int
	Port    string
	RPCPort string
	Db      *env_config.Db
	Cors    *env_config.Cors
}
