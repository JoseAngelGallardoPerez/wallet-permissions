package dao

import (
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/doug-martin/goqu/v9"
)

type Action struct {
	*DAO
}

func (a *Action) NewModel() *model.Action {
	return model.NewAction(a.backend)
}

func (a *Action) Find(handlers ...ConditionHandler) ([]*model.Action, error) {
	proto := a.NewModel()
	rows, err := a.findMany(proto.Model, handlers...)

	if nil != err {
		return nil, err
	}

	result := make([]*model.Action, 0, 10)
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

func (a *Action) FindOne(handlers ...ConditionHandler) (*model.Action, error) {
	Action := a.NewModel()
	row, err := a.findOne(Action.Model, handlers...)
	if nil != err {
		return nil, err
	}

	Action.FillModel(row)
	return Action, nil
}

func (a *Action) Count(handlers ...ConditionHandler) (int64, error) {
	return a.count(model.Tables.Actions(), handlers...)
}

func (a *Action) ById(id int64) (*model.Action, error) {
	return a.FindOne(ByField(model.Tables.ActionsCol("id"), id))
}

func (a *Action) ByName(name string) (*model.Action, error) {
	return a.FindOne(ByField(model.Tables.ActionsCol("name"), name))
}

func (a *Action) ByKey(key string) (*model.Action, error) {
	return a.FindOne(ByField(model.Tables.ActionsCol("key"), key))
}

func (a *Action) ByGroupId(groupId int64) ([]*model.Action, error) {
	return a.Find(func(d *goqu.SelectDataset) *goqu.SelectDataset {
		d = d.InnerJoin(goqu.I(model.Tables.GroupsActions()), goqu.On(goqu.I(model.Tables.ActionsCol("id")).Eq(goqu.I(model.Tables.GroupsActionsCol("action_id")))))
		d = d.Where(goqu.Ex{model.Tables.GroupsActionsCol("group_id"): groupId})
		d = d.Where(goqu.Ex{model.Tables.ActionsCol("is_hidden"): 0})
		return d
	})
}

func (a *Action) IsActionAssignedWithGroup(actionId, groupId int64) (bool, error) {
	action, err := a.FindOne(func(d *goqu.SelectDataset) *goqu.SelectDataset {
		d = d.InnerJoin(goqu.I(model.Tables.GroupsActions()), goqu.On(goqu.I(model.Tables.ActionsCol("id")).Eq(goqu.I(model.Tables.GroupsActionsCol("action_id")))))
		d = d.Where(goqu.Ex{
			model.Tables.GroupsActionsCol("group_id"):  groupId,
			model.Tables.GroupsActionsCol("action_id"): actionId,
		})
		return d
	})

	if nil != err {
		return false, err
	}

	return action.IsExist(), nil
}
