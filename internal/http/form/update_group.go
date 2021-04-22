package form

import "github.com/Confialink/wallet-permissions/internal/types"

type UpdateGroup struct {
	Name        string              `form:"name" json:"name"`
	Description string              `form:"description" json:"description"`
	Scope       string              `form:"scope" json:"scope"`
	Actions     types.ActionsStates `form:"actions" json:"actions"`
}
