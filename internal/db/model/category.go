package model

import (
	"encoding/json"

	"github.com/Confialink/wallet-permissions/internal/db"
)

type Category struct {
	*Model
	name string
	sort int64
}

func NewCategory(backend *db.Backend) *Category {
	model := NewModel(backend, Tables.Categories())
	a := &Category{Model: model}
	model.fieldsValues = []interface{}{
		"id", &a.id,
		"name", &a.name,
		"sort", &a.sort,
	}

	return a
}

func (a *Category) Name() string {
	return a.name
}

func (a *Category) SetName(name string) {
	a.name = name
}

func (a *Category) Sort() int64 {
	return a.sort
}

func (a *Category) SetSort(sort int64) {
	a.sort = sort
}

func (a *Category) Save() error {
	return a.Model.Save()
}

func (a *Category) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"object": "category",
		"id":     a.GetId(),
		"name":   a.Name(),
		"sort":   a.Sort(),
	}

	return json.Marshal(obj)
}
