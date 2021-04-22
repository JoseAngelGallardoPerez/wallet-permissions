package controller

import (
	"net/http"
	"strconv"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/http/response"
	"github.com/Confialink/wallet-permissions/internal/service"
)

//TODO: refactor code duplication
type action struct {
	base
	actionsService    *service.Actions
	permissionService *service.Permissions
}

func (a *action) ListAction(c *gin.Context) {
	var err error
	currentUser := a.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceAction, acl.PrivilegeRead, a.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := a.permissionService.Check(currentUser.UID, service.ActionViewSettings.String())
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

	var actions []*model.Action
	groupId, err := strconv.ParseInt(c.Query("groupId"), 10, 64)

	if groupId != 0 && err == nil {
		actions, err = a.actionsService.FindByGroupId(groupId)
	} else {
		actions, err = a.actionsService.GetAll()
	}

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	// TODO: refactor 999
	resp, err := response.NewList(actions, 999)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.New().SetData(resp))
}

func (a *action) GetAction(c *gin.Context) {
	currentUser := a.mustGetCurrentUser(c)
	isAllowed := acl.ACL.Can(acl.ResourceAction, acl.PrivilegeRead, a.getCurrentUserRole(c))

	if !isAllowed {
		perm, typedErr := a.permissionService.Check(currentUser.UID, service.ActionViewSettings.String())
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

	actionId, typedErr := a.getIdParam(c)
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	action, err := a.actionsService.FindById(actionId)
	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.New().SetData(action))
}
