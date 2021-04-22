package controller

import (
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/http/response"
	"github.com/Confialink/wallet-permissions/internal/service"
)

//TODO: refactor code duplication
type permission struct {
	base
	permissionService *service.Permissions
}

//ListAction lists all permissions with user
//user_id is required
func (p *permission) ListAction(c *gin.Context) {
	currentUser := p.mustGetCurrentUser(c)

	userIdForm := &UserId{}
	if err := c.ShouldBind(userIdForm); err != nil {
		errors.AddShouldBindError(c, err)
		return
	}

	isAllowed, typedErr := p.canRead(currentUser, userIdForm, c)

	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	actions, err := p.permissionService.GetAll(userIdForm.UserId)
	if nil != err {
		privateErr := &errors.PrivateError{Message: "Can't get permissions"}
		privateErr.AddLogPair("error", err.Error())
		errors.AddErrors(c, privateErr)
		return
	}

	resp, err := response.NewWithList(actions, 999)
	if nil != err {
		privateErr := &errors.PrivateError{Message: "Can't form response"}
		privateErr.AddLogPair("error", err.Error())
		errors.AddErrors(c, privateErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (p *permission) GetAction(c *gin.Context) {
	currentUser := p.mustGetCurrentUser(c)

	userIdForm := &UserId{}
	if err := c.ShouldBind(userIdForm); err != nil {
		errors.AddShouldBindError(c, err)
		return
	}

	isAllowed, typedErr := p.canRead(currentUser, userIdForm, c)

	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	if !isAllowed {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	key := c.Param("key")

	perm, typedErr := p.permissionService.Check(userIdForm.UserId, key)
	if nil != typedErr {
		errors.AddErrors(c, typedErr)
		return
	}

	c.JSON(http.StatusOK, response.New().SetData(perm))

}

func (p *permission) canRead(currentUser *userpb.User, userIdForm *UserId, c *gin.Context) (bool, errors.TypedError) {
	//users who are able to edit
	isAllowed := acl.ACL.Can(acl.ResourcePermission, acl.PrivilegeEdit, p.getCurrentUserRole(c))

	//users are able to read own permissions
	if !isAllowed {
		isAllowed = userIdForm.UserId == currentUser.UID
	}

	//admins with "SETTINGS" permission should be able to read permissions of any users
	if !isAllowed {
		perm, typedErr := p.permissionService.Check(currentUser.UID, service.ActionViewSettings.String())
		isAllowed = perm.IsAllowed
		if nil != typedErr {
			errors.AddErrors(c, typedErr)
			return false, typedErr
		}
	}

	return isAllowed, nil
}
