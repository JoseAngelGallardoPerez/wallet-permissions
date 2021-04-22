package util

import (
	"reflect"
	"regexp"
	"strings"
)

type text struct{}

var Text text

func TakeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := TakeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func TakeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

// Converts a string to snake_case
func (*text) ToDelimiter(s, delimiter string) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextIsCapital := false
		if i+1 < len(s) {
			w := s[i+1]
			nextIsCapital = w >= 'A' && w <= 'Z'
		}
		if i > 0 && v >= 'A' && v <= 'Z' && n[len(n)-1] != '_' && !nextIsCapital {
			// add underscore if next letter is a capital
			n += delimiter + string(v)
		} else if v == ' ' {
			// replace spaces with delimiter
			n += delimiter
		} else {
			n = n + string(v)
		}
	}
	n = strings.ToLower(n)
	return n
}

func (t *text) ToSnake(s string) string {
	return t.ToDelimiter(s, "_")
}
