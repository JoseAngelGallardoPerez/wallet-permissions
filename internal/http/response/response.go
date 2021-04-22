package response

import (
	"encoding/json"
	"errors"

	"github.com/Confialink/wallet-permissions/internal/util"
)

type Error struct {
	Title  string `json:"title"`
	Source string `json:"source,omitempty"`
	Detail string `json:"detail,omitempty"`
	Type   string `json:"type"`
	Meta   Meta   `json:"metadata,omitempty"`
}

type Meta interface{}

type List struct {
	HasMore bool        `json:"has_more"`
	Items   interface{} `json:"items"`
}

type Response struct {
	Data     interface{} `json:"data,omitempty"`
	Messages []string    `json:"messages,omitempty"`
	Errors   []*Error    `json:"errors,omitempty"`
}

func New() *Response {
	return new(Response)
}

func NewList(items interface{}, usedLimit uint) (*List, error) {
	values, ok := util.TakeSliceArg(items)
	if !ok {
		return nil, errors.New("items expected to be a slice value")
	}
	result := &List{}
	if usedLimit > 0 && uint(len(values)) > usedLimit {
		result.HasMore = true
		result.Items = values[:len(values)-1]
		return result, nil
	}
	result.HasMore = false
	result.Items = items

	return result, nil
}

func NewWithMessage(message string) *Response {
	return New().AddMessage(message)
}

func NewWithList(items interface{}, usedLimit uint) (*Response, error) {
	list, err := NewList(items, usedLimit)
	if nil != err {
		return nil, err
	}

	res := New()
	res.SetData(list)

	return res, nil
}

func NewWithError(title, errtype string) *Response {
	return New().AddError(title, errtype, nil, nil, nil)
}

func (r *Response) AddError(title, errtype string, source, detail *string, meta Meta) *Response {
	e := &Error{Title: title, Type: errtype}
	if nil != source {
		e.Source = *source
	}
	if nil != detail {

	}
	if nil != meta {
		e.Meta = meta
	}
	r.Errors = append(r.Errors, e)
	return r
}

func (r *Response) AddMessage(message string) *Response {
	r.Messages = append(r.Messages, message)
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"object":   "list",
		"has_more": l.HasMore,
		"items":    l.Items,
	})
}
