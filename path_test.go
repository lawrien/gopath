package gopath

import (
	"testing"
)

type Person struct {
	Name    string
	Age     int
	Friends []Person
}

func TestSimple(t *testing.T) {
	jim := Person{Name: "Jim", Age: 31}

	name := Find("Name", jim)[0].(string)

	if name != jim.Name {
		t.Fail()
	}

	age := Find("Age", jim)[0].(int)
	if age != jim.Age {
		t.Fail()
	}
}

func TestArray(t *testing.T) {
	jim := Person{Name: "Jim", Age: 31}

	jim.Friends = append(jim.Friends, Person{Name: "John", Age: 44})
	jim.Friends = append(jim.Friends, Person{Name: "Claire", Age: 62})

	names := Find("/Friends/*/Name", &jim)

	if len(names) != len(jim.Friends) {
		t.Fail()
	}

	for i, name := range names {
		if jim.Friends[i].Name != name {
			t.Fail()
		}
	}
}

func TestInplaceUpdate(t *testing.T) {
	type Fish struct {
		Name    string
		Spots   map[int]string
		Stripes []int
	}

	f := &Fish{
		Name:    "Harold",
		Spots:   map[int]string{1: "a", 2: "b", 3: "c"},
		Stripes: []int{1, 2, 3, 4, 5},
	}

	m := Find("Spots", f)[0]
	m.(map[int]string)[4] = "d"

	if d, ok := f.Spots[4]; !ok {
		t.Fail()
	} else {
		if d != f.Spots[4] {
			t.Fail()
		}
	}

}
