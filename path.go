package gopath

import (
	"reflect"
	"strings"
)

type Iter struct {
	values []reflect.Value
	index  int
	next   reflect.Value
}

func (self *Iter) Next() bool {
	self.next = reflect.Value{}
	if self.index >= len(self.values) {
		return false
	}
	self.next = self.values[self.index]
	self.index += 1
	return true
}

func (self *Iter) Value() interface{} {
	if self.next.IsValid() {
		return self.next.Interface()
	}
	return nil
}

func isContainer(val reflect.Value) bool {
	if !val.IsValid() {
		return false
	}
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	k := val.Kind()
	return (k == reflect.Struct || k == reflect.Slice || k == reflect.Map)
}

type Path struct {
	path []string
}

func NewPath(path string) *Path {
	if path == "" {
		path = "*"
	} else {
		path = strings.Trim(path, "/")
	}

	splitPath := strings.Split(path, "/")

	p := &Path{path: splitPath}
	return p
}

func (self *Path) Iter(container interface{}) *Iter {
	i := &Iter{}
	i.find(self.path, reflect.ValueOf(container))
	return i
}

func (self *Path) First(container interface{}) (interface{}, bool) {
	i := self.Iter(container)
	if i.Next() {
		return i.Value(), true
	} else {
		return nil, false
	}
}

func (self *Iter) appendValue(val reflect.Value) {
	if !val.IsValid() {
		return
	}
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	self.values = append(self.values, val)
}

func (self *Iter) find(path []string, val reflect.Value) {

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = reflect.Indirect(val)
	}

	if len(path) == 0 {
		self.appendValue(val)
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		self.findStruct(path, val)
	case reflect.Slice:
		self.findSlice(path, val)
	case reflect.Map:
		self.findMap(path, val)
	}
}

func (self *Iter) findStruct(path []string, val reflect.Value) {

	t := val.Type()

	switch {
	case path[0] == "**":
		self.find(path[1:], val)
		for i := 0; i < t.NumField(); i++ {
			value := val.Field(i)

			if isContainer(value) {
				self.find(path, value)
			}
		}

	case path[0] == "*":
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := val.Field(i)
			if field.Anonymous {
				self.find(path, value)
				continue
			}

			self.find(path[1:], value)
		}
	default:
		value := val.FieldByName(path[0])
		if value.IsValid() {
			self.find(path[1:], value)
		}
	}
}

func (self *Iter) findMap(path []string, val reflect.Value) {

	switch {
	case path[0] == "**":
		self.find(path[1:], val)
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			if isContainer(value) {
				self.find(path, value)
			}
		}
	case path[0] == "*":
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			self.find(path[1:], value)
		}
	default:
		value := val.MapIndex(reflect.ValueOf(path[0]))
		if value.IsValid() {
			self.find(path[1:], value)
		}
	}
}

func (self *Iter) findSlice(path []string, val reflect.Value) {

	switch {
	case path[0] == "**":
		self.find(path[1:], val)
		for i := 0; i < val.Len(); i++ {
			value := val.Index(i)

			if isContainer(value) {
				self.find(path, value)
			}
		}
	case path[0] == "*":
		for i := 0; i < val.Len(); i++ {
			value := val.Index(i)
			self.find(path[1:], value)
		}
	}

}
