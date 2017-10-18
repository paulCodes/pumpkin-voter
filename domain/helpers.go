package domain

import (
	"reflect"
	"strings"
	"unicode"
)

func UpcaseToReadable(name string) string {
	var words []string
	l := 0
	for s := name; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, s[:l])
	}
	return strings.Join(words, " ")
}

func GorpTitleDeterminer(e interface{}, fieldName string) string {
	typ := reflect.TypeOf(e)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	field, ok := typ.FieldByName(fieldName)
	if !ok {
		return "GORP TITLE FAILED ON " + fieldName
	}
	tag := field.Tag.Get("gorptitle")
	if tag == "" {
		if fieldName == "Id" {
			return "ID"
		} else {
			return UpcaseToReadable(fieldName)
		}
	}
	return tag
}
