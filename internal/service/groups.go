package service

import (
	"errors"
	"log"

	errorsPackage "github.com/Confialink/wallet-pkg-errors"

	"github.com/Confialink/wallet-permissions/internal/db"
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/types"
)

type Groups struct {
	daoGroup          *dao.Group
	daoAction         *dao.Action
	userService       *User
	db                *db.Backend
	groupsActionsSync *GroupsActionsSynchroniser
}

func NewGroups(
	daoGroup *dao.Group,
	daoAction *dao.Action,
	userService *User,
	db *db.Backend,
	groupsActionsSync *GroupsActionsSynchroniser,
) *Groups {
	return &Groups{daoGroup: daoGroup, daoAction: daoAction, userService: userService, db: db, groupsActionsSync: groupsActionsSync}
}

//CreateGroup creates new group with a given name and associates it with provided actions
func (g *Groups) CreateGroup(name, description, scope string, actionsKeys []string) (group *model.Group, typedError errorsPackage.TypedError) {
	group = g.daoGroup.NewModel()
	group.SetName(name)
	group.SetDescription(description)
	group.SetScope(scope)

	err := group.Save()
	if nil != err {
		privateErr := &errorsPackage.PrivateError{Message: "failed to save new group"}
		privateErr.AddLogPair("error", err.Error())
		typedError = privateErr
		return
	}

	for _, key := range actionsKeys {
		var action *model.Action
		action, err = g.daoAction.ByKey(key)
		if nil != err {
			privateErr := &errorsPackage.PrivateError{Message: "failed to associate action with a group"}
			privateErr.AddLogPair("error", err.Error())
			typedError = privateErr
			return
		}
		allowErr := group.AllowAction(action)
		if nil != allowErr {
			log.Printf("github.com/Confialink/wallet-permissions: unable to assign action \"%s\" - %s", key, err.Error())
		}
	}
	if err := g.groupsActionsSync.SyncPermissionsInheritance(group); err != nil {
		privateErr := &errorsPackage.PrivateError{Message: "failed to sync actions with a group"}
		privateErr.AddLogPair("error", err.Error())
		typedError = privateErr
		return
	}

	return
}

func (g *Groups) DeleteGroup(groupId int64) (typedError errorsPackage.TypedError) {
	admins, err := g.userService.GetByClassId(groupId)
	if err != nil {
		typedError = &errorsPackage.PrivateError{Message: err.Error()}
	}

	if len(admins) > 0 {
		meta := make([]string, len(admins))
		for i, admin := range admins {
			meta[i] = admin.Username
		}
		return errcodes.CreatePublicError(errcodes.CodeAdminGroupConstraint)
	}

	group, err := g.requireGroup(groupId)
	if err != nil {
		typedError = &errorsPackage.PrivateError{Message: err.Error()}
	}

	err = group.Delete()
	if err != nil {
		typedError = &errorsPackage.PrivateError{Message: err.Error()}
	}

	return nil
}

//RenameGroup checks for duplicates and rename passed group
func (g *Groups) RenameGroup(group *model.Group, newName string) errorsPackage.TypedError {
	group.SetName(newName)
	if err := group.Save(); err != nil {
		return &errorsPackage.PrivateError{Message: err.Error()}
	}
	return nil
}

func (g *Groups) RemoveActionFromGroup(actionKey string, groupId int64) error {
	group, err := g.requireGroup(groupId)
	if nil != err {
		return err
	}
	return g.withAction(actionKey, group, "remove")
}

func (g *Groups) FindByName(name string) (*model.Group, errorsPackage.TypedError) {
	group, err := g.daoGroup.ByName(name)
	if err != nil {
		return group, &errorsPackage.PrivateError{Message: err.Error()}
	}
	return group, nil
}

// UpdateActions iterates overs the given actions states and adds or removes actions from group
// according to the "Enable" field
func (g *Groups) UpdateActions(states types.ActionsStates, groupId int64) error {
	var err error
	group, err := g.requireGroup(groupId)
	if nil != err {
		return err
	}
	for _, state := range states {
		if state.Enabled {
			err = g.withAction(state.Key, group, "add")
		} else {
			err = g.withAction(state.Key, group, "remove")
		}
		if err != nil {
			return err
		}
	}

	if err := g.groupsActionsSync.SyncPermissionsInheritance(group); err != nil {
		return err
	}

	return nil
}

//FindByUserId finds group by user id
func (g *Groups) FindByUserId(userId string) (group *model.Group, err error) {
	return g.daoGroup.ByUserId(userId)
}

//FindById finds group by its id
func (g *Groups) FindById(groupId int64) (group *model.Group, err error) {
	return g.daoGroup.ById(groupId)
}

func (g *Groups) GetAll() ([]*model.Group, error) {
	return g.daoGroup.Find()
}

//FindByScope finds group by scope
func (g *Groups) FindByScope(scope string) ([]*model.Group, error) {
	return g.daoGroup.ByScope(scope)
}

//withAction adds action to group or removes action from group
func (g *Groups) withAction(actionKey string, group *model.Group, do string) (err error) {
	action, err := g.daoAction.ByKey(actionKey)
	if nil != err {
		log.Printf("service-permissions: failed to %s action to group - %s", do, err.Error())
	}

	switch do {
	case "add":
		err = group.AllowAction(action)
	case "remove":
		err = group.DenyAction(action)
	default:
		log.Fatalf("invalid argument \"do\" is passed it must be either \"add\" or \"remove\" - \"%s\" is given", do)
	}
	return
}

func (g *Groups) requireGroup(groupId int64) (group *model.Group, err error) {
	group, err = g.FindById(groupId)
	if nil != err {
		return
	}
	if !group.IsExist() {
		log.Printf("github.com/Confialink/wallet-permissions: group with id %d is not found", groupId)
		err = errors.New("group does not exist")
		return
	}
	return
}
