package service

import (
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
	"github.com/Confialink/wallet-permissions/internal/types"
)

type Categories struct {
	daoCategory *dao.Category
	daoAction   *dao.Action
}

func NewCategories(daoCategory *dao.Category, daoAction *dao.Action) *Categories {
	return &Categories{daoCategory: daoCategory, daoAction: daoAction}
}

func (s *Categories) GetAllWithTreePermissions(groupId int64) (result []*types.Category, err error) {
	categories, err := s.daoCategory.Find(dao.OrderDesc("sort"))
	if err != nil {
		return nil, err
	}

	allActions, err := s.daoAction.Find(dao.ByField(model.Tables.ActionsCol("is_hidden"), 0), dao.OrderDesc("sort"))
	if err != nil {
		return nil, err
	}

	selectedActions := make(map[string]*model.Action)
	if groupId > 0 {
		selectedActs, err := s.daoAction.ByGroupId(groupId)
		if err != nil {
			return nil, err
		}
		for _, act := range selectedActs {
			selectedActions[act.Key()] = act
		}
	}

	result = make([]*types.Category, 0, len(categories))
	for _, category := range categories {
		item := &types.Category{
			Id:          category.GetId(),
			Name:        category.Name(),
			Sort:        category.Sort(),
			Permissions: s.permissions(selectedActions, allActions, category),
		}
		result = append(result, item)
	}

	return result, nil
}

// return permissions for category
func (s *Categories) permissions(selectedActions map[string]*model.Action, actions []*model.Action, category *model.Category) []*types.Action {
	result := make([]*types.Action, 0)
	for _, action := range actions {
		if action.CategoryId() != nil && *action.CategoryId() == category.GetId() {
			var enabled bool
			if _, ok := selectedActions[action.Key()]; ok {
				enabled = true
			}
			result = append(result, s.action(selectedActions, actions, action, enabled))
		}
	}

	return result
}

// return children for a permission
func (s *Categories) childPermissions(selectedActions map[string]*model.Action, actions []*model.Action, parentAction *model.Action) []*types.Action {
	result := make([]*types.Action, 0)
	for _, action := range actions {
		if action.ParentId() != nil && *action.ParentId() == parentAction.GetId() {
			var enabled bool
			if _, ok := selectedActions[action.Key()]; ok {
				enabled = true
			}
			result = append(result, s.action(selectedActions, actions, action, enabled))
		}
	}
	return result
}

// create Action struct
func (s *Categories) action(selectedActions map[string]*model.Action, actions []*model.Action, action *model.Action, enabled bool) *types.Action {
	return &types.Action{
		ActionState: types.ActionState{
			Key:     action.Key(),
			Enabled: enabled,
		},
		Id:         action.GetId(),
		ParentId:   action.ParentId(),
		CategoryId: action.CategoryId(),
		Name:       action.Name(),
		Sort:       action.Sort(),
		Children:   s.childPermissions(selectedActions, actions, action),
	}
}
