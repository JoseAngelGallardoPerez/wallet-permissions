package controller

import (
	"strconv"

	"github.com/Confialink/wallet-pkg-errors"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/acl"
	"github.com/Confialink/wallet-permissions/internal/db/model/dao"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
)

type base struct{}

type UserId struct {
	UserId string `form:"userId" json:"userId" binding:"required"`
}

func (*base) getCurrentUser(c *gin.Context) *userpb.User {
	user, exist := c.Get("_user")
	if !exist {
		return nil
	}
	return user.(*userpb.User)
}

func (b *base) getCurrentUserRole(c *gin.Context) acl.Role {
	user := b.getCurrentUser(c)
	if nil == user {
		return acl.Guest
	}
	return acl.RolesHelper.FromName(user.RoleName)
}

func (b *base) mustGetCurrentUser(c *gin.Context) *userpb.User {
	user := b.getCurrentUser(c)
	if nil == user {
		panic("user must be set")
	}
	return user
}

func (*base) int64Param(name string, c *gin.Context) (int64, errors.TypedError) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if nil != err {
		return 0, errcodes.CreatePublicError(errcodes.CodeNumeric)
	}
	return id, nil
}

func filterByField(c *gin.Context, queryParam, columnName string) dao.ConditionHandler {
	return func(d *goqu.SelectDataset) *goqu.SelectDataset {
		values := c.QueryArray(queryParam)
		if len(values) > 0 {
			d = d.Where(goqu.I(columnName).In(values))
			return d
		}
		value := c.Query(queryParam)
		if value != "" {
			d = d.Where(goqu.Ex{columnName: value})
		}
		return d
	}
}

func (b *base) getIdParam(c *gin.Context) (int64, errors.TypedError) {
	return b.int64Param("id", c)
}
