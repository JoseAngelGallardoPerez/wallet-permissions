package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Confialink/wallet-permissions/internal/db"
	"github.com/doug-martin/goqu/v9"
)

type scanner interface {
	Scan(dest ...interface{}) error
}

type modelHydrator interface {
	FillModel(row scanner) error
	GetFields() []string
}

type BaseError struct {
	message string
}

type TypeMismatchErr struct {
	*BaseError
}

type Model struct {
	hydrator modelHydrator
	table    string
	backend  *db.Backend
	// field name must be followed by its value
	// map is not used here b/c of iteration order
	fieldsValues []interface{}
	fields       []string
	values       []interface{}
	fieldsMap    map[string]interface{}
	id           int64

	tx *sql.Tx
}

func (m *Model) Transaction(tx *sql.Tx) {
	m.tx = tx
}

func (m *Model) TransactionOut() {
	m.tx = nil
}

func (m *Model) Table() string {
	return m.table
}

func (e *BaseError) Error() string {
	return e.message
}

func NewModel(backend *db.Backend, table string) *Model {
	model := &Model{backend: backend, table: table}
	model.hydrator = model

	return model
}

func (m *Model) FindById(id int64) (*Model, error) {
	return m.FindOne(goqu.Ex{"id": id})
}

func (m *Model) FindOne(ex ...goqu.Expression) (*Model, error) {
	query, args, err := m.Where(ex...).Limit(1).ToSQL()
	if nil != err {
		return m, err
	}
	row := m.backend.Connection.QueryRow(query, args...)

	m.hydrator.FillModel(row)

	return m, nil
}

func (m *Model) IsExist() bool {
	return m.GetId() != 0
}

func (m *Model) GetId() int64 {
	return m.id
}

func (m *Model) FieldExist(name string) (ok bool) {
	m.initFieldAndValues()
	_, ok = m.fieldsMap[name]
	return
}

//SetField sets value to the corresponding field by it name
//though it might be useful, try to avoid using this method
//if possible since it works much slower than direct set methods
//TODO: refactor - create private setters like (setStringField, setIntField etc.), lazy map child fields with setters
func (m *Model) SetField(name string, value interface{}) (ok bool, err error) {
	m.initFieldAndValues()
	ref, ok := m.fieldsMap[name]
	if !ok {
		err = errors.New("Field " + name + " does not exist")
		return
	}

	errMsg := "Cannot set field value: passed value is type of %s, but field \"" + name + "\" has a different type"
	switch value.(type) {
	case string:
		var r *string
		r, ok = ref.(*string)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "string"))
			return
		}
		*r = value.(string)
	case int:
		var r *int
		r, ok = ref.(*int)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "int"))
			return
		}
		*r = value.(int)
	case int32:
		var r *int32
		r, ok = ref.(*int32)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "int32"))
			return
		}
		*r = value.(int32)
	case int64:
		var r *int64
		r, ok = ref.(*int64)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "int64"))
			return
		}
		*r = value.(int64)
	case float32:
		var r *float32
		r, ok = ref.(*float32)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "float32"))
			return
		}
		*r = value.(float32)
	case float64:
		var r *float64
		r, ok = ref.(*float64)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "float64"))
			return
		}
		*r = value.(float64)
	case bool:
		var r *bool
		r, ok = ref.(*bool)
		if !ok {
			err = typeMismatchError(fmt.Sprintf(errMsg, "bool"))
			return
		}
		*r = value.(bool)
	default:
		panic("Invalid type specified")

	}

	return
}

func (m *Model) Delete() error {

	if !m.IsExist() {
		return nil
	}

	query := fmt.Sprintf("DELETE from `%s` WHERE id = %d LIMIT 1", m.table, m.id)

	_, err := m.exec(query)

	if nil != err {
		log.Println("Error: deleting record ", err)
	}

	m.id = 0

	return err
}

func (m *Model) GetFields() []string {
	m.initFieldAndValues()
	return m.fields
}

func (m *Model) FillModel(row scanner) error {
	err := row.Scan(append([]interface{}{&m.id}, m.getValues()...)...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (m *Model) getValues() []interface{} {
	m.initFieldAndValues()
	return m.values
}

func (m *Model) initFieldAndValues() {
	if len(m.fields) == 0 {
		m.fields = make([]string, len(m.fieldsValues)/2)
		m.values = make([]interface{}, len(m.fieldsValues)/2)
		m.fieldsMap = make(map[string]interface{}, len(m.fieldsValues)/2)
		fieldIndex := 0
		valueIndex := 0
		isField := true
		field := ""
		for _, value := range m.fieldsValues {
			if isField {
				field = value.(string)
				m.fields[fieldIndex] = field
				fieldIndex++
				isField = false
				continue
			}
			m.values[valueIndex] = value
			m.fieldsMap[field] = value
			valueIndex++
			isField = true
		}
	}
}

func (m *Model) setId(id int64) *Model {
	m.id = id
	return m
}

func (m *Model) Where(expressions ...goqu.Expression) *goqu.SelectDataset {
	return m.SelectFrom().Where(expressions...).Prepared(true)
}

func (m *Model) SelectFrom() *goqu.SelectDataset {
	selectFields := make([]interface{}, len(m.GetFields())+1)
	selectFields[0] = m.table + ".id"

	i := 1
	for _, v := range m.GetFields() {
		selectFields[i] = m.table + "." + v
		i++
	}

	qb := m.backend.Builder.From(m.table).Select(selectFields...)

	return qb
}

func (m Model) copy() *Model {
	return &m
}

func (m *Model) Save() error {
	if m.IsExist() {
		return update(m)
	}
	return insert(m)
}

func insert(m *Model) error {
	var fields = make([]string, len(m.GetFields()))
	var values = make([]interface{}, len(m.getValues()))
	copy(fields, m.GetFields())
	copy(values, m.getValues())

	for i := 0; i < len(fields); i++ {
		if fields[i] == "updated_at" {
			fields = append(fields[:i], fields[i+1:]...)
			values = append(values[:i], values[i+1:]...)
			break
		}
	}

	query := insertTemplate(m.table, fields...)

	stmt, err := m.prepare(query)

	if nil != err {
		log.Println("Error: preparing query\""+query+"\" ", err)
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(values...)

	if nil != err {
		log.Println("Error: executing query\""+query+"\" ", err)
		return err
	}

	m.id, err = result.LastInsertId()

	if nil != err {
		log.Println("Error: ", err)
		return err
	}

	return err
}

func fastBulkInsert(models []*Model) error {

	var model *Model
	rowsNumber := len(models)

	if rowsNumber == 0 {
		return nil
	}

	model = models[0]

	valuesNumber := rowsNumber * len(model.GetFields())

	for _, f := range model.GetFields() {
		if f == "id" {
			valuesNumber = rowsNumber * (len(model.GetFields()) - 1)
		}
	}

	query := bulkInsertTemplate(model.table, rowsNumber, model.GetFields()...)

	values := make([]interface{}, valuesNumber)

	i := 0
	for _, m := range models {
		for _, v := range m.getValues() {
			values[i] = v
			i++
		}
	}

	stmt, err := model.prepare(query)

	if nil != err {
		log.Println("Error: preparing query\""+query+"\" ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(values...)

	if nil != err {
		log.Println("Error: executing query\""+query+"\" ", err)
		return err
	}

	return err
}

func update(m *Model) error {
	var fields = make([]string, len(m.GetFields()))
	var values = make([]interface{}, len(m.getValues())+1)
	copy(fields, m.GetFields())
	copy(values, append(m.getValues(), m.id))

	for i := 0; i < len(fields); i++ {
		if fields[i] == "created_at" {
			fields = append(fields[:i], fields[i+1:]...)
			values = append(values[:i], values[i+1:]...)
			break
		}
	}

	query := updateTemplate(m.table, fields...)

	stmt, err := m.prepare(query)

	if nil != err {
		log.Println("Error: preparing query\""+query+"\" ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(values...)

	if nil != err {
		log.Println("Error: executing query\""+query+"\" ", err)
		return err
	}

	return err
}

func (m *Model) exec(query string, args ...interface{}) (sql.Result, error) {
	if nil != m.tx {
		return m.tx.Exec(query, args...)
	}
	return m.backend.Connection.Exec(query, args...)
}

func (m *Model) prepare(query string) (*sql.Stmt, error) {
	if nil != m.tx {
		return m.tx.Prepare(query)
	}
	return m.backend.Connection.Prepare(query)
}

func typeMismatchError(message string) *TypeMismatchErr {
	return &TypeMismatchErr{&BaseError{message: message}}
}
