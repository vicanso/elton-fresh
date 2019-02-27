# cod-fresh

[![Build Status](https://img.shields.io/travis/vicanso/cod-fresh.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-fresh)

HTTP response freshness testing middleware for cod.

```go
package main

import (
	"bytes"

	"github.com/vicanso/cod"

	etag "github.com/vicanso/cod-etag"
	fresh "github.com/vicanso/cod-fresh"
)

func main() {

	d := cod.New()
	d.Use(fresh.NewDefault())
	d.Use(etag.NewDefault())

	d.GET("/", func(c *cod.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("abcd")
		return
	})

	d.ListenAndServe(":7001")
}
```