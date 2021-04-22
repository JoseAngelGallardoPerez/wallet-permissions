package controller

import "github.com/Confialink/wallet-permissions/internal/app/di"

var c = di.Container
var Factory *factory

type factory struct{}

func init() {
	Factory = &factory{}
}

func (*factory) PermissionController() (ctrl *permission) {
	return &permission{
		permissionService: c.ServicePermissions(),
	}
}

func (*factory) GroupController() (ctrl *group) {
	return &group{
		permissionService: c.ServicePermissions(),
		groupsService:     c.ServiceGroups(),
	}
}

func (*factory) ActionController() (ctrl *action) {
	return &action{
		permissionService: c.ServicePermissions(),
		actionsService:    c.ServiceActions(),
	}
}

func (*factory) CategoryController() (ctrl *category) {
	return &category{
		permissionService: c.ServicePermissions(),
		service:           c.ServiceCategories(),
	}
}
