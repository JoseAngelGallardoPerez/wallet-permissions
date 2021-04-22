package service

import (
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
)

type Action string

const (
	ActionViewSettings   = Action("view_settings")
	ActionModifySettings = Action("modify_settings")
	ActionCreateSettings = Action("create_settings")
	ActionRemoveSettings = Action("remove_settings")
)

func (a Action) String() string {
	return string(a)
}

type Actions struct {
	daoAction *dao.Action
}

func NewActions(daoAction *dao.Action) *Actions {
	return &Actions{daoAction: daoAction}
}

func (a *Actions) FindById(actionId int64) (*model.Action, error) {
	return a.daoAction.ById(actionId)
}

func (a *Actions) FindByGroupId(groupId int64) ([]*model.Action, error) {
	return a.daoAction.ByGroupId(groupId)
}

func (a *Actions) GetAll() ([]*model.Action, error) {
	return a.daoAction.Find()
}
