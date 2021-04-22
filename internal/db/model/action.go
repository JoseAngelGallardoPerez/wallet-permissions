package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Confialink/wallet-permissions/internal/db"
)

type Action struct {
	*Model
	name        string
	key         string
	description sql.NullString
	createdAt   time.Time
	updatedAt   *time.Time
	categoryId  *int64
	parentId    *int64
	isHidden    bool // some actions are not used at current time but frontend checks them
	sort        int64
}

func NewAction(backend *db.Backend) *Action {
	model := NewModel(backend, Tables.Actions())
	a := &Action{Model: model}
	model.fieldsValues = []interface{}{
		"name", &a.name,
		"key", &a.key,
		"description", &a.description,
		"created_at", &a.createdAt,
		"updated_at", &a.updatedAt,
		"category_id", &a.categoryId,
		"parent_id", &a.parentId,
		"is_hidden", &a.isHidden,
		"sort", &a.sort,
	}

	return a
}

func (a *Action) UpdatedAt() *time.Time {
	return a.updatedAt
}

func (a *Action) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Action) Description() sql.NullString {
	return a.description
}

func (a *Action) SetDescription(description sql.NullString) {
	a.description = description
}

func (a *Action) Name() string {
	return a.name
}

func (a *Action) SetName(name string) {
	a.name = name
}

func (a *Action) CategoryId() *int64 {
	return a.categoryId
}

func (a *Action) SetCategoryId(categoryId *int64) {
	a.categoryId = categoryId
}

func (a *Action) IsHidden() bool {
	return a.isHidden
}

func (a *Action) SetIsHidden(isHidden bool) {
	a.isHidden = isHidden
}

func (a *Action) Sort() int64 {
	return a.sort
}

func (a *Action) SetSort(sort int64) {
	a.sort = sort
}

func (a *Action) ParentId() *int64 {
	return a.parentId
}

func (a *Action) SetParentId(parentId *int64) {
	a.parentId = parentId
}

func (a *Action) Key() string {
	return a.key
}

func (a *Action) SetKey(key string) {
	a.key = key
}

func (a *Action) Save() error {
	a.createdAt = time.Now()
	now := time.Now()
	a.updatedAt = &now

	return a.Model.Save()
}

func (a *Action) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"object":      "action",
		"id":          a.GetId(),
		"name":        a.Name(),
		"key":         a.Key(),
		"description": a.Description().String,
		"createdAt":   a.CreatedAt().Unix(),
		"updatedAt":   nil,
		"categoryId":  a.CategoryId(),
		"parentId":    a.ParentId(),
		"sort":        a.Sort(),
	}
	if nil != a.UpdatedAt() {
		obj["updatedAt"] = a.UpdatedAt().Unix()
	}
	return json.Marshal(obj)
}
