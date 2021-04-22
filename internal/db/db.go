package db

import (
	"database/sql"
	"fmt"

	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/doug-martin/goqu/v9"

	// dialects work like drivers in go where they are not registered until you import the package
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Backend struct {
	Connection *sql.DB
	Builder    *goqu.Database
}

func NewBackend(cfg *env_config.Db) (*Backend, error) {
	db, err := sql.Open(cfg.Driver, fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Schema))

	if nil != err {
		return nil, err
	}

	err = db.Ping()
	if nil != err {
		db.Close()
		return nil, err
	}

	return &Backend{Connection: db, Builder: goqu.New(cfg.Driver, db)}, nil
}
