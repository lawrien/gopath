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

	find(splitPath, container, &results)

	return results
}

func find(path []string, container interface{}, results *[]interface{}) {
	c := reflect.ValueOf(container)
	if c.Kind() == reflect.Ptr {
		if c.IsNil() {
			return
		} else {
			c = reflect.Indirect(c)
		}
	}

	if len(path) == 0 {
		appendValue(results, container)
		return
	}

	switch c.Kind() {
	case reflect.Struct:
		findStruct(path, container, results)
	case reflect.Slice:
		findSlice(path, container, results)
	case reflect.Map:
		findMap(path, container, results)
	}
}

func appendValue(results *[]interface{}, value interface{}) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Struct:
		nv := reflect.New(val.Type()).Elem()
		nv.Set(val)
		val = nv
	}
	// Always put pointers to value where possible
	if val.CanAddr() {
		val = val.Addr()
	}
	*results = append(*results, val.Interface())
}

func findStruct(path []string, container interface{}, results *[]interface{}) {
	c := reflect.ValueOf(container)
	if c.Kind() == reflect.Ptr {
		c = reflect.Indirect(c)
	}
	t := c.Type()

	switch {
	case path[0] == "**":
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := c.Field(i)
			if field.Anonymous {
				find(path, value.Interface(), results)
				continue
			}
			if path[1] == field.Name {
				find(path[2:], value.Interface(), results)
			} else {
				find(path, value.Interface(), results)
			}
		}
	case path[0] == "*":
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := c.Field(i)
			if field.Anonymous {
				find(path, value.Interface(), results)
				continue
			}

			find(path[1:], value.Interface(), results)
		}
	default:
		value := c.FieldByName(path[0])
		if value.IsValid() {
			find(path[1:], value.Interface(), results)
		}
	}
}

func findMap(path []string, container interface{}, results *[]interface{}) {
	c := reflect.ValueOf(container)
	if c.Kind() == reflect.Ptr {
		c = reflect.Indirect(c)
	}

	switch {
	case path[0] == "**":
		for _, key := range c.MapKeys() {
			value := c.MapIndex(key)
			if path[1] == key.String() {
				find(path[2:], value.Interface(), results)
			} else {
				find(path, value.Interface(), results)
			}
		}
	case path[0] == "*":
		for _, key := range c.MapKeys() {
			value := c.MapIndex(key)
			find(path[1:], value.Interface(), results)
		}
	default:
		value := c.MapIndex(reflect.ValueOf(path[0]))
		if value.IsValid() {
			find(path[1:], value.Interface(), results)
		}
	}
}

func findSlice(path []string, container interface{}, results *[]interface{}) {
	c := reflect.ValueOf(container)
	if c.Kind() == reflect.Ptr {
		c = reflect.Indirect(c)
	}

	switch {
	case path[0] == "**":
		for i := 0; i < c.Len(); i++ {
			value := c.Index(i)
			find(path, value.Interface(), results)
		}
	case path[0] == "*":
		for i := 0; i < c.Len(); i++ {
			value := c.Index(i)
			find(path[1:], value.Interface(), results)
		}
	}

}
