package acl

import (
	base "github.com/kildevaeld/go-acl"
)

type acl struct {
	list *base.ACL
}

var (
	list     = base.New(base.NewMemoryStore())
	ACL  acl = acl{list}
)

func init() {
	list.Role(Guest.String(), "")
	list.Role(Client.String(), Guest.String())
	list.Role(Buyer.String(), Client.String())
	list.Role(Supplier.String(), Buyer.String())
	list.Role(Financier.String(), Supplier.String())
	list.Role(Admin.String(), Financier.String())
	list.Role(Root.String(), Admin.String())

	list.Allow(Root, "*", "*")
	list.Allow(Client, PrivilegeRead.String(), ResourcePermission)
	list.Allow(Buyer, PrivilegeRead.String(), ResourcePermission)
	list.Allow(Supplier, PrivilegeRead.String(), ResourcePermission)
	list.Allow(Financier, PrivilegeRead.String(), ResourcePermission)

	list.Allow(Buyer, PrivilegeRead.String(), ResourceGroup)
	list.Allow(Supplier, PrivilegeRead.String(), ResourceGroup)
	list.Allow(Financier, PrivilegeRead.String(), ResourceGroup)
	list.Allow(Admin, PrivilegeRead.String(), ResourceGroup)
}

func (a *acl) Can(resource interface{}, privilege Privilege, roles ...Role) bool {
	for _, role := range roles {
		if a.list.Can(role, "*", "*") {
			return true
		}

		if a.list.Can(role, privilege.String(), "*") {
			return true
		}

		if a.list.Can(role, privilege.String(), resource) {
			return true
		}
	}
	return false
}
