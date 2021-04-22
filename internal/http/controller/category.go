package controller

import (
	"net/http"
	"strconv"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/http/response"
	"github.com/Confialink/wallet-permissions/internal/service"
	"github.com/Confialink/wallet-permissions/internal/types"
)

type category struct {
	base
	service           *service.Categories
	permissionService *service.Permissions
}

func (a *category) List(c *gin.Context) {
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

	groupId, _ := strconv.ParseInt(c.Query("groupId"), 10, 64)
	categories, err := a.service.GetAllWithTreePermissions(groupId)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	// todo: move to the event bus when card module will be as an extension
	categoriesShorten := make([]*types.Category, 0, len(categories))
	cardModuleSettings, err := service.CardModuleSettings()
	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	for _, cat := range categories {
		if !cardModuleSettings.IsEnabled && cat.Name == "Cards" {
			continue
		}
		categoriesShorten = append(categoriesShorten, cat)
	}

	resp, err := response.NewList(categoriesShorten, 0)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.New().SetData(resp))
}
