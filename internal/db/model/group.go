package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Confialink/wallet-permissions/internal/db"
	"github.com/doug-martin/goqu/v9"
)

type Group struct {
	*Model
	name        string
	description sql.NullString
	scope       string
	createdAt   time.Time
	updatedAt   *time.Time
}

func NewGroup(backend *db.Backend) *Group {
	model := NewModel(backend, Tables.Groups())
	g := &Group{Model: model}
	model.fieldsValues = []interface{}{
		"name", &g.name,
		"scope", &g.scope,
		"description", &g.description,
		"created_at", &g.createdAt,
		"updated_at", &g.updatedAt,
	}

	return g
}

//AddUser adds user to a group
func (g *Group) AddUser(userId string) error {
	insert := g.backend.Builder.From(Tables.UsersGroups()).Insert().Rows(goqu.Record{
		Tables.UsersGroupsCol("user_id"):  userId,
		Tables.UsersGroupsCol("group_id"): g.GetId(),
	})

	_, err := insert.Executor().Exec()
	return err
}

//AllowAction marks action as allowed for a group
func (g *Group) AllowAction(action *Action) error {
	if !action.IsExist() {
		return errors.New("action does not exist")
	}
	insert := g.backend.Builder.From(Tables.GroupsActions()).Insert().Rows(goqu.Record{
		Tables.GroupsActionsCol("group_id"):  g.GetId(),
		Tables.GroupsActionsCol("action_id"): action.GetId(),
	}).OnConflict(goqu.DoNothing())

	_, err := insert.Executor().Exec()
	return err
}

//DenyAction deletes association with an action
func (g *Group) DenyAction(action *Action) error {
	if !action.IsExist() {
		return errors.New("action does not exist")
	}
	del := g.backend.Builder.From(Tables.GroupsActions()).Where(goqu.Ex{
		Tables.GroupsActionsCol("action_id"): action.GetId(),
		Tables.GroupsActionsCol("group_id"):  g.GetId(),
	}).Delete()

	_, err := del.Executor().Exec()
	return err
}

func (g *Group) UpdatedAt() *time.Time {
	return g.updatedAt
}

func (g *Group) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) SetName(name string) {
	g.name = name
}

func (g *Group) Scope() string {
	return g.scope
}

func (g *Group) SetScope(scope string) {
	g.scope = scope
}

func (g *Group) Description() string {
	return g.description.String
}

func (g *Group) SetDescription(description string) {
	g.description.String = description
	g.description.Valid = true
}

func (g *Group) Save() error {
	g.createdAt = time.Now()
	now := time.Now()
	g.updatedAt = &now

	return g.Model.Save()
}

func (g *Group) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"object":      "group",
		"id":          g.GetId(),
		"name":        g.Name(),
		"scope":       g.Scope(),
		"description": g.Description(),
		"createdAt":   g.CreatedAt().Unix(),
		"updatedAt":   nil,
	}
	if nil != g.UpdatedAt() {
		obj["updatedAt"] = g.UpdatedAt().Unix()
	}
	return json.Marshal(obj)
}
