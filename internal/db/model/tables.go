package model

var tables = map[string]string{
	"categories":     "categories",
	"actions":        "actions",
	"groups":         "groups",
	"groups_actions": "groups_actions",
	"users_groups":   "users_groups",
}

type tablesHelper struct{}

var Tables = tablesHelper{}

func (t *tablesHelper) Categories() string {
	return t.tableName("categories")
}

func (t *tablesHelper) CategoriesCol(name string) string {
	return columnName(t.tableName("categories"), name)
}

func (t *tablesHelper) Actions() string {
	return t.tableName("actions")
}

func (t *tablesHelper) ActionsCol(name string) string {
	return columnName(t.tableName("actions"), name)
}

func (t *tablesHelper) Groups() string {
	return t.tableName("groups")
}

func (t *tablesHelper) GroupsCol(name string) string {
	return columnName(t.tableName("groups"), name)
}

func (t *tablesHelper) GroupsActions() string {
	return t.tableName("groups_actions")
}

func (t *tablesHelper) GroupsActionsCol(name string) string {
	return columnName(t.tableName("groups_actions"), name)
}

func (t *tablesHelper) UsersGroups() string {
	return t.tableName("users_groups")
}

func (t *tablesHelper) UsersGroupsCol(name string) string {
	return columnName(t.tableName("users_groups"), name)
}

func columnName(table, column string) string {
	return table + "." + column
}

func (t *tablesHelper) tableName(name string) string {
	return "permissions_" + tables[name]
}
