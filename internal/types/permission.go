package types

import "encoding/json"

type Permission struct {
	ActionKey string
	UserId    string
	IsAllowed bool
}

func (p *Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"object":    "permission",
		"actionKey": p.ActionKey,
		"userId":    p.UserId,
		"isAllowed": p.IsAllowed,
	})
}
