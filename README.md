# gopath

Simple mechanism to allow the use of `xpath` style paths to access data in Go `map`s, `struct`s and slices.
Currently early days, so dont expect too much, but it does fundamentally work.

## Usage

Import the package using the `go get` command.

    go get github.com/lawrien/gopath

## API
Finding items is done by creating in three steps:

  * Create a Path object using the `gopath.NewPath` API

        path := gopath.NewPath("*/People/Ages")

  * Create and iterator for the path, passing in the object to be searched:

        iter := path.Iter(myobj)    

  * Use the Next() and Value() APIs of the Iterator to walk through the available results:

        for iter.Next() {
          v := iter.Value()
        }    

### Path Syntax

Nothing fancy here, just three different parts available to the path:

  - `**` - path wildcard. Will match all possible paths.
  - `*` - path element wildcard. Will match a single part of the path
  - `key or fieldname` - matches a fieldname or a map key name

### Examples

      import "github.com/lawrien/gopath"
      jim := Person{Name: "Jim", Age: 31}

      jim.Friends = append(jim.Friends, Person{Name: "John", Age: 44})
      jim.Friends = append(jim.Friends, Person{Name: "Claire", Age: 62})
      
      it := NewPath("/Friends/*/Name").Iter(jim)
      for i := 0; it.Next(); i++ {
        name := it.Value().(string)
        fmt.Printf("Friend -> %s\n", name)
      }
    

### Limitations/Be aware

OK, this is early code, so could do with more testing, examples and docs. So, before using,
here's some of the limitations at the moment.

Firstly, paths are not checked for syntax ! It would be very easy to put in a dud path. Here's
some simple rules:

  - No need to start paths with a '/' - it's assumed
  - Dont have more than one `**` adjacent in your paths (ie `**/**/fish`). It will probably 
    blow up.
  - No method for accessing individual array elements at the moment.

As with all things Go, `gopath` can only operate on `public` structs and fields. Anything 
private will not be found.

If you need a pointer to an item (ie so you can modify the contents), use the `Iter.ValuePtr()` API rather than `Iter.Value()`.

Dont be surprised that accessing a slice returns a slice ! If you want the content of a slice, 
you should use a path like `myarray/*`.



