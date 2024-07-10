# go-where

[![Go Reference](https://pkg.go.dev/badge/github.com/coreybutler/go-where.svg)](https://pkg.go.dev/github.com/coreybutler/go-where)

A library for determining the root path of an executable. Supported on Windows, macOS, and Linux.

_myapp.go_:

```go
package main

import (
  "fmt"
  "os"
  "github.com/coreybutler/go-where"
)

func main() {
  executable := os.Args[1]
  path, err := where.Find(executable)

  if err != nil {
    panic(err)
  }

  fmt.Print(path)
}
```

Run this with:

```sh
$ go run myapp.go node.exe
C:\nodejs\node.exe
```

## Exceptions

There are environments where the same executable may exist on a system multiple times (different versions, different architecture, etc). If you need to exclude known specific paths, use the `FindExcept`.

```go
package main

import (
  "fmt"
  "os"
  "github.com/coreybutler/go-where"
)

func main() {
  executable := os.Args[1]
  // The boolean argument indicates a recursive lookup
  path, err := where.FindExcept(executable, true, "C:\nodejs\node.exe")

  if err != nil {
    panic(err)
  }

  fmt.Print(path)
}
```

```sh
$ go run myapp.go node.exe
C:\nvm\v20.0.0\node.exe
# Notice C:\nodejs\node.exe was ignored
```

## Alternative function

`where.FindAll()` and `where.FindAllExcept()` are also available. These methods return a slice of strings (`[]string`) containing the paths where the executable/binary is located (except those explicitly ignored in `FindAllExcept`). This is useful for identifying multiple installations of a particular program. An empty slice is returned if the file cannot be found.

## File Extensions

It is best to supply the file extension of the executable, but this library will attempt to identify executables in two manners.

**By Permissions**
If the file has explicit execute permissions, it will be considered "executable" and returned by the `Find/FindExcept` and `FindAll/FindAllExcept` methods.

**By extension**
This module attempts to determine if a file is executable by its file extension.

For example:

```go
fmt.Print(where.Find("node"))

// C:\nodejs\node.exe
```

File extension identification is limited to a hard coded list of known extensions. See the files named `expand_*.go`.

--

Copyright (c) 2021 Corey Butler and contributors.
