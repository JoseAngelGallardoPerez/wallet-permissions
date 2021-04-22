module github.com/Confialink/wallet-permissions

go 1.13

replace github.com/Confialink/wallet-permissions/rpc/permissions => ./rpc/permissions

require (
	github.com/Confialink/wallet-permissions/rpc/permissions v0.0.0-00010101000000-000000000000
	github.com/Confialink/wallet-pkg-discovery/v2 v2.0.0-20210217105157-30e31661c1d1
	github.com/Confialink/wallet-pkg-env_config v0.0.0-20210217112253-9483d21626ce
	github.com/Confialink/wallet-pkg-env_mods v0.0.0-20210217112432-4bda6de1ee2c
	github.com/Confialink/wallet-pkg-errors v1.0.2
	github.com/Confialink/wallet-pkg-service_names v0.0.0-20210217112604-179d69540dea
	github.com/Confialink/wallet-settings/rpc/proto/settings v0.0.0-20210218070334-b4153fc126a0
	github.com/Confialink/wallet-users/rpc/proto/users v0.0.0-20210218071418-0600c0533fb2
	github.com/doug-martin/goqu/v9 v9.9.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/inconshreveable/log15 v0.0.0-20200109203555-b30bc20e4fd1
	github.com/kildevaeld/go-acl v0.0.0-20171228130000-7799b11f4759
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/twitchtv/twirp v5.12.0+incompatible
)
