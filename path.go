package gopath

import (
	"reflect"
	"strings"
)

func Find(path string, container interface{}) []interface{} {
	var results = []interface{}{}
	if path == "" {
		path = "*"
	} else {
		path = strings.Trim(path, "/")
	}

	splitPath := strings.Split(path, "/")

	find(splitPath, reflect.ValueOf(container), &results)

	return results
}

func FindClean(path string, container interface{}) []interface{} {
	var results = Find(path, container)
	var clean = []interface{}{}
	var m = map[interface{}]string{}

	for _, val := range results {
		m[val] = ""
	}

	for key, _ := range m {
		if reflect.ValueOf(key) != reflect.Zero(reflect.TypeOf(key)) {
			clean = append(clean, key)
		}
	}

	return clean
}

func find(path []string, val reflect.Value, results *[]interface{}) {
	pv := val
	if pv.Kind() == reflect.Ptr {
		if pv.IsNil() {
			return
		}
		val = pv.Elem()
	}

	if len(path) == 0 {
		appendValue(results, pv)
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		findStruct(path, pv, results)
	case reflect.Slice:
		findSlice(path, pv, results)
	case reflect.Map:
		findMap(path, pv, results)
	}
}

func appendValue(results *[]interface{}, val reflect.Value) {
	pv := val
	if pv.Kind() == reflect.Ptr {
		val = reflect.Indirect(pv)
	}

	if val.Kind() == reflect.Struct {
		if val.CanAddr() {
			*results = append(*results, val.Addr().Interface())
			return
		}
	}
	*results = append(*results, val.Interface())
}

func findStruct(path []string, val reflect.Value, results *[]interface{}) {
	pv := val
	if pv.Kind() == reflect.Ptr {
		if pv.IsNil() {
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		val = pv.Elem()
	}

	t := val.Type()

	switch {
	case path[0] == "**":
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := val.Field(i)
			if field.Anonymous {
				find(path, value, results)
				continue
			}
			if path[1] == field.Name {
				find(path[2:], value, results)
			} else {
				find(path, value, results)
			}
		}
	case path[0] == "*":
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := val.Field(i)
			if field.Anonymous {
				find(path, value, results)
				continue
			}

			find(path[1:], value, results)
		}
	default:
		value := val.FieldByName(path[0])
		if value.IsValid() {
			find(path[1:], value, results)
		}
	}
}

func findMap(path []string, val reflect.Value, results *[]interface{}) {
	pv := val
	if pv.Kind() == reflect.Ptr {
		if pv.IsNil() {
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		val = pv.Elem()
	}

	switch {
	case path[0] == "**":
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			if path[1] == key.String() {
				find(path[2:], value, results)
			} else {
				find(path, value, results)
			}
		}
	case path[0] == "*":
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			find(path[1:], value, results)
		}
	default:
		value := val.MapIndex(reflect.ValueOf(path[0]))
		if value.IsValid() {
			find(path[1:], value, results)
		}
	}
}

func findSlice(path []string, val reflect.Value, results *[]interface{}) {
	pv := val
	if pv.Kind() == reflect.Ptr {
		if pv.IsNil() {
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		val = pv.Elem()
	}

	switch {
	case path[0] == "**":
		for i := 0; i < val.Len(); i++ {
			value := val.Index(i)
			find(path, value, results)
		}
	case path[0] == "*":
		for i := 0; i < val.Len(); i++ {
			value := val.Index(i)
			find(path[1:], value, results)
		}
	}

}
