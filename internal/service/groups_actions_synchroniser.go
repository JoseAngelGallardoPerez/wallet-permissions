package service

import (
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
)

type GroupsActionsSynchroniser struct {
	daoAction *dao.Action
}

func NewGroupsActionsSynchroniser(
	daoAction *dao.Action,
) *GroupsActionsSynchroniser {
	return &GroupsActionsSynchroniser{daoAction}
}

// if an action has a parent we must add group permissions for all parents
func (s *GroupsActionsSynchroniser) SyncPermissionsInheritance(group *model.Group) error {
	actions, err := s.daoAction.ByGroupId(group.GetId())
	if err != nil {
		return err
	}

	for _, action := range actions {
		if err := s.addPermissionGroupActionIfNeed(actions, action, group); err != nil {
			return err
		}
	}

	return nil
}

// add new Actions to the group if Action has a parent Action.
// the method checks all parents recursive
func (s *GroupsActionsSynchroniser) addPermissionGroupActionIfNeed(actions []*model.Action, action *model.Action, group *model.Group) error {
	if action.ParentId() == nil {
		return nil
	}

	parentAction, err := s.daoAction.ById(*action.ParentId())
	if err != nil {
		return err
	}

	// check if parent permission exists for the group
	for _, act := range actions {
		if act.Key() == parentAction.Key() {
			return nil
		}
	}

	// insert parent permission for the group
	if err := group.AllowAction(parentAction); err != nil {
		return err
	}

	// check parent for parent Action
	if err := s.addPermissionGroupActionIfNeed(actions, parentAction, group); err != nil {
		return err
	}

	return nil
}
