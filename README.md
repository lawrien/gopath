# gopath

Simple mechanism to allow the use of `xpath` style paths to access data in Go `map`s, `struct`s and slices.
Currently early days, so dont expect too much, but it does fundamentally work.

## Usage

Import the package using the `go get` command.

    go get github.com/lawrien/gopath

## API
There is a single public method that will accept a path and container (any valid `struct`,`map` or `slice`
at this stage). 

    gopath.Find(path string, container interface{}) []interface{}

Any data that matches the given path will be returned in an array. 

### Path Syntax

Nothing fancy here, just three different parts available to the path:

  - `**` - path wildcard. Will match all possible paths.
  - `*` - path element wildcard. Will match a single part of the path
  - `key or fieldname` - matches a fieldname or a map key name

### Examples

      import "github.com/lawrien/gopath"
      jim := Person{Name: "Jim", Age: 31}

      name := gopath.Find("Name", jim)[0].(string)
      age := gopath.Find("Age", jim)[0].(int)

      // Give Jim some friends
      jim.Friends = append(jim.Friends, Person{Name: "John", Age: 44})
      jim.Friends = append(jim.Friends, Person{Name: "Claire", Age: 62})

      for i, name := range Find("Friends/Name", jim) {
        fmt.Printf("Jim has %s as a friend.\n",name)
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

Currently, a `struct`s are returned as pointers to the struct 
data, not as a copy - this does allow in-place updating of data, but be aware as type
assertations might not work the way you think:

    x := gopath.Find("/path/to/mystruct",someobj)[0]
    x.(*Mystruct) // right
    x.(Mystruct)  // wrong

Dont be surprised that accessing a slice returns a slice ! If you want the content of a slice, 
you should use a path like `gopath.Find("/myarray/*",someobj)`.



