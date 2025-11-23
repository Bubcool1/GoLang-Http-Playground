package main

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

func FirstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

func handleConversion(v any) any {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil
	}
	typ := val.Type()

	// Handle pointer
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		list := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			list[i] = handleConversion(val.Index(i).Interface())
		}
		return list
	}

	if typ.Kind() == reflect.Struct {
		result := make(map[string]any)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			key := FirstToLower(field.Name)
			result[key] = val.Field(i).Interface()
		}
		return result
	}

	return v
}

func StructToMap(v any) any {
	return handleConversion(v)
}
