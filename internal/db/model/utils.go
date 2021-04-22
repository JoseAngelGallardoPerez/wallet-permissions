package model

import (
	"fmt"
	"strings"
)

func repeatPlaceholders(repeat int, fieldsLen int) string {
	tpl := strings.Trim(strings.Repeat("?,", fieldsLen), ",")
	return strings.Trim(strings.Repeat("("+tpl+"),", repeat), ",")
}

func insertTemplate(table string, fields ...string) string {
	placeholders := ""
	fieldsTpl := ""
	for _, f := range fields {
		if f == "id" {
			continue
		}
		placeholders += "?,"
		fieldsTpl += "`" + table + "`." + "`" + f + "`,"
	}
	placeholders = strings.Trim(placeholders, ",")
	fieldsTpl = strings.Trim(fieldsTpl, ",")

	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, fieldsTpl, placeholders)
}

func bulkInsertTemplate(table string, rowsNumber int, fields ...string) string {
	fieldsTpl := ""
	valuesNumber := len(fields)
	for _, f := range fields {
		if f == "id" {
			valuesNumber--
			continue
		}

		fieldsTpl += "`" + table + "`." + "`" + f + "`,"
	}

	fieldsTpl = strings.Trim(fieldsTpl, ",")
	placeholders := repeatPlaceholders(rowsNumber, valuesNumber)

	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s", table, fieldsTpl, placeholders)
}

func updateTemplate(table string, fields ...string) string {
	setTpl := ""
	for _, f := range fields {
		if f == "id" {
			continue
		}
		setTpl += "`" + table + "`." + "`" + f + "` = ?,"

	}
	setTpl = strings.Trim(setTpl, ",")

	return fmt.Sprintf("UPDATE `%s` SET %s WHERE `%s`.`id` = ?", table, setTpl, table)
}
