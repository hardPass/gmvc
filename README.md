# gmvc

web framework for golang

## Installation

To install gmvc:

    go get github.com/hujh/gmvc

## Quick Start
    
```go
package main

import (
	"log"
	"net/http"
	"github.com/hujh/gmvc"
)

func main() {
	app := gmvc.NewApp()

	app.HandleFunc("/", func(c *gmvc.Context) error {
		return c.WriteString("hello world!")
	})

	log.Fatal(s.ListenAndServe(":8080", app))
}
```
