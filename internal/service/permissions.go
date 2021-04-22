package service

import (
	"log"

	errorsPackage "github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/doug-martin/goqu/v9"

	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/types"
)

type Permissions struct {
	daoGroup    *dao.Group
	daoAction   *dao.Action
	userService *User
}

func NewPermissions(daoGroup *dao.Group, daoAction *dao.Action, userService *User) *Permissions {
	return &Permissions{
		daoGroup:    daoGroup,
		daoAction:   daoAction,
		userService: userService,
	}
}

//Check verifies if a given user has a permission
func (p *Permissions) Check(userId, actionKey string) (permission *types.Permission, typedError errorsPackage.TypedError) {
	action, err := p.daoAction.ByKey(actionKey)
	if nil != err {
		privateErr := &errorsPackage.PrivateError{Message: "unable to check permission"}
		privateErr.AddLogPair("error", err)
		typedError = privateErr
		return
	}
	permission = &types.Permission{ActionKey: actionKey, UserId: userId}
	if !action.IsExist() {
		typedError = errcodes.CreatePublicError(errcodes.CodeActionKeyDoesNotExist)
		return
	}

	user, err := p.getUser(userId)
	if err != nil {
		return nil, &errorsPackage.PrivateError{Message: err.Error()}
	}

	if acl.RolesHelper.FromName(user.RoleName).Is(acl.Root) {
		permission.IsAllowed = true
		return permission, nil
	}

	group, err := p.getUserGroup(user)
	if nil != err {
		privateErr := &errorsPackage.PrivateError{Message: "unable to check permission"}
		privateErr.AddLogPair("error", err)
		typedError = privateErr
		return
	}

	if !group.IsExist() {
		permission.IsAllowed = false
		return
	}

	isAssigned, err := p.daoAction.IsActionAssignedWithGroup(action.GetId(), group.GetId())
	if nil != err {
		privateErr := &errorsPackage.PrivateError{Message: "unable to check permission"}
		privateErr.AddLogPair("error", err)
		typedError = privateErr
		return
	}
	permission.IsAllowed = isAssigned
	return
}

func (p *Permissions) Can(userId, actionKey string) (bool, errorsPackage.TypedError) {
	perm, typedErr := p.Check(userId, actionKey)
	if nil != typedErr {
		return false, typedErr
	}

	return perm.IsAllowed, nil
}

func (p *Permissions) CanAny(userId string, actionKeys ...string) (bool, errorsPackage.TypedError) {
	for _, actionKey := range actionKeys {
		can, typedErr := p.Can(userId, actionKey)
		if nil != typedErr {
			return false, typedErr
		}
		if can {
			return can, nil
		}
	}
	return false, nil
}

func (p *Permissions) CheckAll(userId string, actionKeys []string) (result []*types.Permission, typedErr errorsPackage.TypedError) {
	if len(actionKeys) > 100 {
		typedErr = &errorsPackage.PrivateError{Message: "too many actions are given"}
		return
	}
	result = make([]*types.Permission, 0, len(actionKeys))
	for _, actionKey := range actionKeys {
		permission, typedErr := p.Check(userId, actionKey)
		if nil != typedErr {
			return result, typedErr
		}
		result = append(result, permission)
	}

	return
}

func (p *Permissions) GetAll(userId string) (result []*types.Permission, err error) {
	result = make([]*types.Permission, 0, 32)

	user, err := p.getUser(userId)
	if err != nil {
		return nil, err
	}

	group, err := p.getUserGroup(user)

	if nil != err {
		log.Println("github.com/Confialink/wallet-permissions: unable to find group", err)
		return
	}

	if !group.IsExist() {
		return
	}
	userAssignedActions, err := p.daoAction.ByGroupId(group.GetId())
	if nil != err {
		log.Println("github.com/Confialink/wallet-permissions: unable to retrieve permissions", err)
		return
	}
	allActions, err := p.GetAllNotHiddenActions()
	permissionsMap := make(map[string]*types.Permission)

	// if user role is "root" then everything is allowed by default
	defaultAllowed := acl.RolesHelper.FromName(user.RoleName).Is(acl.Root)
	for _, action := range allActions {
		permissionsMap[action.Key()] = &types.Permission{
			ActionKey: action.Key(),
			UserId:    userId,
			IsAllowed: defaultAllowed,
		}
		result = append(result, permissionsMap[action.Key()])
	}
	//assigned userAssignedActions
	for _, action := range userAssignedActions {
		permissionsMap[action.Key()].IsAllowed = true
	}

	return
}

func (p *Permissions) GetAllActions() ([]*model.Action, error) {
	return p.daoAction.Find()
}

func (p *Permissions) GetAllNotHiddenActions() ([]*model.Action, error) {
	return p.daoAction.Find(dao.ByField(model.Tables.ActionsCol("is_hidden"), 0))
}

// GetGroupsByIds returns groups by passed ids
func (p *Permissions) GetGroupsByIds(ids []int64) ([]*model.Group, error) {
	return p.daoGroup.Find(p.groupByIdsCondition(ids))
}

func (p *Permissions) groupByIdsCondition(ids []int64) func(dataSet *goqu.SelectDataset) *goqu.SelectDataset {
	return func(dataSet *goqu.SelectDataset) *goqu.SelectDataset {
		return dataSet.Where(goqu.Ex{
			"id": goqu.Op{"in": ids},
		})
	}
}

func (p *Permissions) getUserGroup(user *users.User) (*model.Group, error) {
	group, err := p.daoGroup.ById(user.AdministratorClassId)
	if nil != err {
		return nil, err
	}

	return group, nil
}

func (p *Permissions) getUser(uid string) (*users.User, error) {
	user, err := p.userService.GetByUID(uid)
	if nil != err {
		log.Println("github.com/Confialink/wallet-permissions: unable to retrieve user", err)
		return nil, err
	}
	return user, nil
}
