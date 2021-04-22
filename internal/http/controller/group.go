package controller

import (
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/http/response"
	"github.com/Confialink/wallet-permissions/internal/service"
	"github.com/Confialink/wallet-permissions/internal/types"
	"github.com/Confialink/wallet-permissions/internal/validator"
)

//TODO: refactor code duplication
type group struct {
	base
	permissionService *service.Permissions
	groupsService     *service.Groups
}

func (g *group) ListAction(c *gin.Context) {
	var groups []*model.Group
	var err error

	isAllowed := acl.ACL.Can(acl.ResourceGroup, acl.PrivilegeRead, g.getCurrentUserRole(c))

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	userId := c.Query("userId")
	scope := c.Query("scope")

	if userId != "" {
		foundGroup, err := g.groupsService.FindByUserId(userId)
		if nil != err {
			errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
			return
		}
		if foundGroup.IsExist() {
			groups = append(make([]*model.Group, 0, 1), foundGroup)
		}
	} else {
		if scope != "" {
			groups, err = g.groupsService.FindByScope(scope)
		} else {
			groups, err = g.groupsService.GetAll()
		}
	}

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	// TODO: refactor 999
	resp, err := response.NewList(groups, 999)
	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.New().SetData(resp))
}

func (g *group) GetAction(c *gin.Context) {
	currentUser := g.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceGroup, acl.PrivilegeRead, g.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := g.permissionService.Check(currentUser.UID, service.ActionViewSettings.String())
		isAllowed = perm.IsAllowed
		if nil != typedErr {
			errors.AddErrors(c, typedErr)
			return
		}
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	groupId, typedErr := g.getIdParam(c)

	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	foundGroup, err := g.groupsService.FindById(groupId)
	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.New().SetData(foundGroup))
}

func (g *group) CreateAction(c *gin.Context) {
	currentUser := g.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceGroup, acl.PrivilegeEdit, g.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := g.permissionService.Check(currentUser.UID, service.ActionCreateSettings.String())
		isAllowed = perm.IsAllowed
		if nil != typedErr {
			errors.AddErrors(c, typedErr)
			return
		}
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	groupForm := &struct {
		Name          string              `form:"name" json:"name" binding:"required,groupNameUniqueness"`
		ActionsStates types.ActionsStates `form:"actions" json:"actions"`
		Description   string              `form:"description" json:"description"`
		Scope         string              `form:"scope" json:"scope" binding:"required"`
	}{}

	if err := c.ShouldBind(groupForm); err != nil {
		errors.AddShouldBindError(c, err)
		return
	}

	groupModel, typedErr := g.groupsService.CreateGroup(groupForm.Name, groupForm.Description, groupForm.Scope, groupForm.ActionsStates.ActiveKeys())
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	c.JSON(http.StatusOK, groupModel)
}

func (g *group) RemoveActionAction(c *gin.Context) {
	currentUser := g.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceGroup, acl.PrivilegeEdit, g.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := g.permissionService.Check(currentUser.UID, service.ActionRemoveSettings.String())
		isAllowed = perm.IsAllowed
		if nil != typedErr {
			errors.AddErrors(c, typedErr)
			return
		}
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	groupId, typedErr := g.getIdParam(c)
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	groupModel, err := g.groupsService.FindById(groupId)
	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	if !groupModel.IsExist() {
		errcodes.AddError(c, errcodes.CodeGroupDoesNotExist)
		return
	}

	actionKey := c.Param("action")

	g.groupsService.RemoveActionFromGroup(actionKey, groupId)

	c.JSON(http.StatusOK, response.NewWithMessage("successfully deleted"))
}

func (g *group) UpdateAction(c *gin.Context) {
	currentUser := g.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceGroup, acl.PrivilegeEdit, g.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := g.permissionService.Check(currentUser.UID, service.ActionModifySettings.String())
		isAllowed = perm.IsAllowed
		if nil != typedErr {
			errors.AddErrors(c, typedErr)
			return
		}
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	groupId, typedErr := g.getIdParam(c)
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	groupModel, err := g.groupsService.FindById(groupId)
	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	if !groupModel.IsExist() {
		errcodes.AddError(c, errcodes.CodeGroupDoesNotExist)
		return
	}

	groupForm, typedErr := validator.UpdateGroup(c, groupModel, g.groupsService)
	if typedErr != nil {
		errors.AddErrors(c, typedErr)
		return
	}

	if groupForm.Description != "" {
		groupModel.SetDescription(groupForm.Description)
		groupModel.Save()
	}

	if groupForm.Name != "" && groupForm.Name != groupModel.Name() {
		typedErr = g.groupsService.RenameGroup(groupModel, groupForm.Name)
		if nil != err {
			errors.AddErrors(c, typedErr)
			return
		}
	}

	if len(groupForm.Actions) > 0 {
		if err = g.groupsService.UpdateActions(groupForm.Actions, groupId); err != nil {
			errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, response.New().SetData(groupModel))
}

func (g *group) DeleteAction(c *gin.Context) {
	currentUser := g.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceGroup, acl.PrivilegeEdit, g.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := g.permissionService.Check(currentUser.UID, service.ActionRemoveSettings.String())
		isAllowed = perm.IsAllowed
		if nil != typedErr {
			errors.AddErrors(c, typedErr)
			return
		}
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	groupId, typedErr := g.getIdParam(c)
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	typedErr = g.groupsService.DeleteGroup(groupId)
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	c.JSON(http.StatusOK, response.NewWithMessage("successfully deleted"))
}
