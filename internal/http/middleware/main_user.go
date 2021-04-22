package middleware

import (
	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
)

// MainUser is user with role_name buyer, supplier or financier than can to create staff users
func MainUserOnly(c *gin.Context) {
	user, exist := c.Get("_user")
	if !exist {
		errcodes.AddError(c, errcodes.CodeForbidden)
		c.Abort()
		return
	}

	role := acl.RolesHelper.FromName((user.(*userpb.User)).RoleName)
	if role < acl.Buyer {
		errcodes.AddError(c, errcodes.CodeForbidden)
		c.Abort()
		return
	}
}
