# Setting Up a Go Project

## Creating the root project directory
```bash
mkdir go-project
cd go-project
```

## Initialise the project as a Go module
```bash
go mod init source/go-project-name
```
If the module is designed to be published such that the Go package manager
can download it, `source` must refer to the URL of where the module is
hosted.

For more information, refer to:
https://go.dev/doc/modules/managing-dependencies#naming_module

## Craete the `main` package for your project
```go
package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, World!")
}
```
