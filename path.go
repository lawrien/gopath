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

	// if val.Kind() == reflect.Struct {
	// 	if val.CanAddr() {
	// 		*results = append(*results, val.Addr().Interface())
	// 		return
	// 	}
	// }
	self.values = append(self.values, val)
}

func (self *Iter) find(path []string, val reflect.Value) {
	pv := val
	if pv.Kind() == reflect.Ptr {
		if pv.IsNil() {
			return
		}
		val = pv.Elem()
	}

	if len(path) == 0 {
		self.appendValue(pv)
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		self.findStruct(path, pv)
	case reflect.Slice:
		self.findSlice(path, pv)
	case reflect.Map:
		self.findMap(path, pv)
	}
}

func (self *Iter) findStruct(path []string, val reflect.Value) {
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
				self.find(path, value)
				continue
			}
			if path[1] == field.Name {
				self.find(path[2:], value)
			} else {
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
				self.find(path[2:], value)
			} else {
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
			self.find(path, value)
		}
	case path[0] == "*":
		for i := 0; i < val.Len(); i++ {
			value := val.Index(i)
			self.find(path[1:], value)
		}
	}

}
