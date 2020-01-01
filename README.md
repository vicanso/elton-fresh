# elton-fresh

[![Build Status](https://img.shields.io/travis/vicanso/elton-fresh.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-fresh)

HTTP response freshness testing middleware for elton.

```go
package main

import (
	"bytes"

	"github.com/vicanso/elton"

	etag "github.com/vicanso/elton-etag"
	fresh "github.com/vicanso/elton-fresh"
)

func main() {

	e := elton.New()
	e.Use(fresh.NewDefault())
	e.Use(etag.NewDefault())

	e.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("abcd")
		return
	})

	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
```