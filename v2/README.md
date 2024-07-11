# go-where

[![Go Reference](https://pkg.go.dev/badge/github.com/coreybutler/go-where.svg)](https://pkg.go.dev/github.com/coreybutler/go-where)

A library for determining the root path of an executable. Supported on Windows, macOS, and Linux.

> **Breaking changes in v2.0.0**
> The library has been simplified to a single exported function that accepts a configuration object. All of the prior capabilities are supported, but there are no longer inidivdual functions for finding one path vs multiple paths, exceptions, etc.

_myapp.go_:

```go
package main

import (
  "fmt"
  "os"
  "github.com/coreybutler/go-where/v2"
)

func main() {
  executable := os.Args[1]
  path, err := where.Find(executable)

  if err != nil {
    panic(err)
  }

  fmt.Print(path[0])
}
```

Run this with:

```sh
$ go run myapp.go node.exe
C:\nodejs\node.exe
```

## Confguration Options

The `Options` type is defined as:

```go
type Options struct {
	All       bool     `json:"all"`
	Recursive bool     `json:"recursive"`
	Except    []string `json:"except"`
}
```

| Option        | Description                                                                                                                           |           Default           |
| ------------- | :------------------------------------------------------------------------------------------------------------------------------------ | :--------------------------: |
| _All_       | Return all paths where the executable is found (as opposed to the first one)                                                          |          `false`          |
| _Recursive_ | Search `PATH` directories recursively for the executable.                                                                           |           `true`           |
| _Except_    | A slice of paths/glob patterns to ignore. Ths can be used to override specific file types, such as ignoring `.bat` files on Windows | `[]string{}` (empty slice) |

```go
package main

import (
  "fmt"
  "os"
  "github.com/coreybutler/go-where/v2"
)

func main() {
  executable := os.Args[1]
  // The boolean argument indicates a recursive lookup
  path, err := where.Find(executable, where.Options{
    Except: []string{"C:\nodejs\node.exe"}
  })

  if err != nil {
    panic(err)
  }

  fmt.Print(path[0])
}
```

```sh
$ go run myapp.go node.exe
not found
# C:\nodejs\node.exe was ignored!
```

## File Extensions

It is best to supply the file extension of the executable, but this library will attempt to identify executables in three ways:

**By bit**
On some operating systems, the first byte(s) of the file flag whether it is executable or not.

**By file permissions**
If the file has explicit execute permissions, it will be considered "executable".

**By file extension**
File extensions are used as a last resort.

For example:

```go
fmt.Print(where.Find("node"))
// C:\nodejs\node.exe
```

File extension identification is limited to a hard-coded list of known extensions. See the files named `expand_*.go`.

---

Copyright (c) 2021-2024 Corey Butler and contributors.
