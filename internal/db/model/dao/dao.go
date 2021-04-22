package dao

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/Confialink/wallet-permissions/internal/db"
	"github.com/Confialink/wallet-permissions/internal/db/model"
)

type DAO struct {
	backend *db.Backend
}

type ConditionHandler func(d *goqu.SelectDataset) *goqu.SelectDataset

func New(backend *db.Backend) *DAO {
	return &DAO{backend: backend}
}

func ByField(field string, value interface{}) ConditionHandler {
	return func(d *goqu.SelectDataset) *goqu.SelectDataset {
		return d.Where(goqu.Ex{field: value})
	}
}

func OrderDesc(field string) ConditionHandler {
	return func(d *goqu.SelectDataset) *goqu.SelectDataset {
		return d.Order(goqu.I(field).Desc())
	}
}

func (d *DAO) findMany(proto *model.Model, handlers ...ConditionHandler) (*sql.Rows, error) {
	qb := proto.SelectFrom()
	for _, h := range handlers {
		qb = h(qb)
	}
	query, args, err := qb.ToSQL()

	if nil != err {
		return nil, err
	}
	stmt, err := d.backend.Connection.Prepare(query)
	if nil != err {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if nil != err {
		return nil, err
	}

	return rows, nil
}

func (d *DAO) count(table string, handlers ...ConditionHandler) (int64, error) {
	qb := d.backend.Builder.From(table).Prepared(true)
	for _, h := range handlers {
		qb = h(qb)
	}

	return qb.Count()
}

func (d *DAO) findOne(proto *model.Model, handlers ...ConditionHandler) (*sql.Row, error) {
	qb := proto.SelectFrom()
	for _, h := range handlers {
		qb = h(qb)
	}
	query, args, err := qb.ToSQL()

	if nil != err {
		return nil, err
	}

	return d.backend.Connection.QueryRow(query, args...), nil
}
