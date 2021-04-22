package dao

import (
	"github.com/Confialink/wallet-permissions/internal/db/model"
)

type Category struct {
	*DAO
}

func (a *Category) NewModel() *model.Category {
	return model.NewCategory(a.backend)
}

func (a *Category) Find(handlers ...ConditionHandler) ([]*model.Category, error) {
	proto := a.NewModel()
	rows, err := a.findMany(proto.Model, handlers...)

	if nil != err {
		return nil, err
	}

	result := make([]*model.Category, 0, 10)
	for rows.Next() {
		item := a.NewModel()
		err = item.FillModel(rows)
		if nil != err {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
