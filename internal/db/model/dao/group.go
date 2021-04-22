package dao

import (
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/doug-martin/goqu/v9"
)

type Group struct {
	*DAO
}

func (g *Group) NewModel() *model.Group {
	return model.NewGroup(g.backend)
}

func (g *Group) Find(handlers ...ConditionHandler) ([]*model.Group, error) {
	proto := g.NewModel()
	rows, err := g.findMany(proto.Model, handlers...)

	if nil != err {
		return nil, err
	}

	result := make([]*model.Group, 0, 10)
	for rows.Next() {
		item := g.NewModel()
		err = item.FillModel(rows)
		if nil != err {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (g *Group) FindOne(handlers ...ConditionHandler) (*model.Group, error) {
	Group := g.NewModel()
	row, err := g.findOne(Group.Model, handlers...)
	if nil != err {
		return nil, err
	}

	Group.FillModel(row)
	return Group, nil
}

func (g *Group) Count(handlers ...ConditionHandler) (int64, error) {
	return g.count(model.Tables.Groups(), handlers...)
}

func (g *Group) ById(id int64) (*model.Group, error) {
	return g.FindOne(ByField(model.Tables.GroupsCol("id"), id))
}

func (g *Group) ByName(name string) (*model.Group, error) {
	return g.FindOne(ByField(model.Tables.GroupsCol("name"), name))
}

func (g *Group) ByScope(scope string) ([]*model.Group, error) {
	return g.Find(ByField(model.Tables.GroupsCol("scope"), scope))
}

func (g *Group) ByUserId(userId string) (*model.Group, error) {
	return g.FindOne(func(d *goqu.SelectDataset) *goqu.SelectDataset {
		d = d.InnerJoin(goqu.I(model.Tables.UsersGroups()), goqu.On(goqu.I(model.Tables.GroupsCol("id")).Eq(goqu.I(model.Tables.UsersGroupsCol("group_id")))))
		d = d.Where(goqu.Ex{model.Tables.UsersGroupsCol("user_id"): userId})
		return d
	})
}
