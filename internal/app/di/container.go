package di

import (
	"os"
	"strconv"

	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/Confialink/wallet-pkg-service_names"
	"github.com/go-playground/validator/v10"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-permissions/internal/app"
	"github.com/Confialink/wallet-permissions/internal/app/config"
	"github.com/Confialink/wallet-permissions/internal/db"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
	"github.com/Confialink/wallet-permissions/internal/rpc"
	"github.com/Confialink/wallet-permissions/internal/service"
)

var Container *container

type DataAccessObject struct {
	dao      *dao.DAO
	action   *dao.Action
	category *dao.Category
	group    *dao.Group
}

type container struct {
	appConfig config.Main
	dbBackend *db.Backend

	rpcServer *rpc.Server

	services struct {
		groups            *service.Groups
		permissions       *service.Permissions
		actions           *service.Actions
		categories        *service.Categories
		user              *service.User
		groupsActionsSync *service.GroupsActionsSynchroniser
	}

	DataAccessObject DataAccessObject

	validate *validator.Validate
	logger   log15.Logger
}

func init() {
	Container = newContainer()
	readConfig(&Container.appConfig, Container.ServiceLogger().New("service", "configReader"))
	Container.Validator()
}

func newContainer() *container {
	result := &container{}
	return result
}

func (c *container) Config() *config.Main {
	return &c.appConfig
}

func (c *container) DbBackend() *db.Backend {
	var err error
	if nil == c.dbBackend {
		c.dbBackend, err = db.NewBackend(c.appConfig.Db)
		if nil != err {
			c.ServiceLogger().Error("Can't establish connection to backend", "error", err)
			panic(err)
		}
	}

	return c.dbBackend
}

func (c *container) RPCServer() *rpc.Server {
	if nil == c.rpcServer {
		c.rpcServer = rpc.NewRPCServer(c.ServicePermissions())
	}
	return c.rpcServer
}

func (c *container) ServiceGroups() *service.Groups {
	if nil == c.services.groups {
		d := c.DataAccessObject
		c.services.groups = service.NewGroups(d.Group(), d.Action(), c.ServiceUser(), c.DbBackend(), c.ServiceGroupsActionsSync())
	}
	return c.services.groups
}

func (c *container) ServiceGroupsActionsSync() *service.GroupsActionsSynchroniser {
	if nil == c.services.groupsActionsSync {
		d := c.DataAccessObject
		c.services.groupsActionsSync = service.NewGroupsActionsSynchroniser(d.Action())
	}
	return c.services.groupsActionsSync
}

func (c *container) ServicePermissions() *service.Permissions {
	if nil == c.services.permissions {
		d := c.DataAccessObject
		c.services.permissions = service.NewPermissions(d.Group(), d.Action(), c.ServiceUser())
	}
	return c.services.permissions
}

func (c *container) ServiceActions() *service.Actions {
	if nil == c.services.actions {
		d := c.DataAccessObject
		c.services.actions = service.NewActions(d.Action())
	}
	return c.services.actions
}

func (c *container) ServiceCategories() *service.Categories {
	if nil == c.services.categories {
		d := c.DataAccessObject
		c.services.categories = service.NewCategories(d.Category(), d.Action())
	}
	return c.services.categories
}

func (c *container) ServiceUser() *service.User {
	if nil == c.services.user {
		c.services.user = service.NewUser()
	}
	return c.services.user
}

func (c *container) ServiceLogger() log15.Logger {
	if c.logger == nil {
		c.logger = log15.New("microservice", service_names.Permissions.Internal)
	}
	return c.logger
}

func (c *container) Validator() *validator.Validate {
	if c.validate == nil {
		c.validate = app.LoadValidator(c.ServiceGroups(), c.ServiceLogger().New("service", "validator"))
	}
	return c.validate
}

func (d *DataAccessObject) Get() *dao.DAO {
	if nil == d.dao {
		d.dao = dao.New(Container.DbBackend())
	}
	return d.dao
}

func (d *DataAccessObject) Action() *dao.Action {
	if nil == d.action {
		d.action = &dao.Action{d.Get()}
	}
	return d.action
}

func (d *DataAccessObject) Category() *dao.Category {
	if nil == d.category {
		d.category = &dao.Category{d.Get()}
	}
	return d.category
}

func (d *DataAccessObject) Group() *dao.Group {
	if nil == d.group {
		d.group = &dao.Group{d.Get()}
	}
	return d.group
}

func readConfig(cfg *config.Main, logger log15.Logger) {
	defaultConfigReader := env_config.NewReader("permissions")
	cfg.Cors = defaultConfigReader.ReadCorsConfig()
	cfg.Db = defaultConfigReader.ReadDbConfig()
	cfg.Env = env_config.Env("ENV", env_mods.Development)
	threads := os.Getenv("VELMIE_WALLET_PERMISSIONS_THREADS")
	if threads != "" {
		iThreads, err := strconv.ParseInt(threads, 10, 32)
		if nil != err {
			Container.ServiceLogger().Error("error reading environment variable \"VELMIE_WALLET_PERMISSIONS_THREADS\"", "error", err)
		} else {
			cfg.Threads = int(iThreads)
		}
	}

	cfg.Port = env_config.Env("VELMIE_WALLET_PERMISSIONS_PORT", "")
	cfg.RPCPort = env_config.Env("VELMIE_WALLET_PERMISSIONS_RPC_PORT", "")

	validateConfig(cfg, logger)
}

func validateConfig(cfg *config.Main, logger log15.Logger) {
	vdtr := env_config.NewValidator(logger)
	vdtr.ValidateCors(cfg.Cors, logger)
	vdtr.ValidateDb(cfg.Db, logger)
	vdtr.CriticalIfEmpty(cfg.Port, "VELMIE_WALLET_PERMISSIONS_PORT", logger)
	vdtr.CriticalIfEmpty(cfg.RPCPort, "VELMIE_WALLET_PERMISSIONS_RPC_PORT", logger)
}
