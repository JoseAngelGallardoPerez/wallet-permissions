package types

type ActionState struct {
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

type ActionsStates []*ActionState

func (a ActionsStates) ActiveKeys() []string {
	result := make([]string, 0, len(a))
	for _, action := range a {
		if action.Enabled {
			result = append(result, action.Key)
		}
	}
	return result
}
