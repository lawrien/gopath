package gopath

import (
	"fmt"
	"testing"
)

type Person struct {
	Name    string
	Age     int
	Friends []Person
}

func TestSimple(t *testing.T) {
	jim := Person{Name: "Jim", Age: 31}

	if name, ok := NewPath("Name", &jim).First(); ok {
		if name.(string) != jim.Name {
			t.Fail()
		}
		fmt.Printf("Name => %s\n", name)
	} else {
		t.Fail()
	}

	if age, ok := NewPath("Age", &jim).First(); ok {
		if age.(int) != jim.Age {
			t.Fail()
		}
		fmt.Printf("Age => %d\n", age)
	} else {
		t.Fail()
	}
}

// func TestArray(t *testing.T) {
// 	jim := Person{Name: "Jim", Age: 31}

// 	jim.Friends = append(jim.Friends, Person{Name: "John", Age: 44})
// 	jim.Friends = append(jim.Friends, Person{Name: "Claire", Age: 62})

// 	names := Find("/Friends/*/Name", &jim)

// 	if len(names) != len(jim.Friends) {
// 		t.Fail()
// 	}

// 	for i, name := range names {
// 		if jim.Friends[i].Name != name {
// 			t.Fail()
// 		}
// 	}
// }

// func TestSlice(t *testing.T) {
// 	s := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}

// 	d := Find("*", s)[0].([]int)
// 	for i, n := range d {
// 		fmt.Printf("%d = %d\n", i, n)
// 	}
// 	*(&d) = append(d, 13)

// 	for i, n := range s[0] {
// 		fmt.Printf("%d = %d\n", i, n)
// 	}
// }

// func TestInplaceUpdate(t *testing.T) {
// 	type Fish struct {
// 		Name    string
// 		Spots   map[int]string
// 		Stripes []int
// 	}

// 	f := Fish{
// 		Name:    "Harold",
// 		Spots:   map[int]string{1: "a", 2: "b", 3: "c"},
// 		Stripes: []int{1, 2, 3, 4, 5},
// 	}

// 	m := Find("Spots", f)[0]
// 	m.(map[int]string)[4] = "d"

// 	if d, ok := f.Spots[4]; !ok {
// 		t.Fail()
// 	} else {
// 		if d != f.Spots[4] {
// 			t.Fail()
// 		}
// 	}
// }
