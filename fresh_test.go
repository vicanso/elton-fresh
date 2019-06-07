package fresh

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestFresh(t *testing.T) {
	fn := NewDefault()
	modifiedAt := "Tue, 25 Dec 2018 00:02:22 GMT"
	t.Run("skip", func(t *testing.T) {
		assert := assert.New(t)
		c := cod.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn := New(Config{
			Skipper: func(c *cod.Context) bool {
				return true
			},
		})
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
	})

	t.Run("return error", func(t *testing.T) {
		assert := assert.New(t)
		c := cod.NewContext(nil, nil)
		customErr := errors.New("abccd")
		c.Next = func() error {
			return customErr
		}
		fn := New(Config{})
		err := fn(c)
		assert.Equal(err, customErr, "custom error should be return")
	})

	t.Run("not modified", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(cod.HeaderIfModifiedSince, modifiedAt)
		resp := httptest.NewRecorder()
		resp.Header().Set(cod.HeaderLastModified, modifiedAt)

		c := cod.NewContext(resp, req)
		done := false
		c.Next = func() error {
			done = true
			c.Body = map[string]string{
				"name": "tree.xie",
			}
			c.BodyBuffer = bytes.NewBufferString(`{"name":"tree.xie"}`)
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.True(done)

		assert.Equal(c.StatusCode, 304, "status code should be 304")
		assert.Nil(c.Body, "body should be nil")
		assert.Nil(c.BodyBuffer, "body buffer should be nil")
	})

	t.Run("no body", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(cod.HeaderIfModifiedSince, modifiedAt)
		resp := httptest.NewRecorder()
		resp.Header().Set(cod.HeaderLastModified, modifiedAt)
		c := cod.NewContext(resp, req)
		c.Next = func() error {
			return nil
		}
		c.NoContent()
		err := fn(c)
		assert.Nil(err)
		assert.Equal(c.StatusCode, 204, "no body should be passed by fresh")
	})

	t.Run("post method", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("POST", "/users/me", nil)
		req.Header.Set(cod.HeaderIfModifiedSince, modifiedAt)
		resp := httptest.NewRecorder()
		resp.Header().Set(cod.HeaderLastModified, modifiedAt)

		c := cod.NewContext(resp, req)
		done := false
		c.Next = func() error {
			done = true
			c.StatusCode = http.StatusOK
			c.Body = map[string]string{
				"name": "tree.xie",
			}
			c.BodyBuffer = bytes.NewBufferString(`{"name":"tree.xie"}`)
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.True(done)

		assert.Equal(c.StatusCode, 200, "post requset should be passed by fresh")
		assert.NotNil(c.Body, "post requset should be passed by fresh")
		assert.NotNil(c.BodyBuffer, "post requset should be passed by fresh")
	})

	t.Run("error response", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(cod.HeaderIfModifiedSince, modifiedAt)
		resp := httptest.NewRecorder()
		resp.Header().Set(cod.HeaderLastModified, modifiedAt)

		c := cod.NewContext(resp, req)
		done := false
		c.Next = func() error {
			done = true
			c.StatusCode = http.StatusBadRequest
			c.Body = map[string]string{
				"name": "tree.xie",
			}
			c.BodyBuffer = bytes.NewBufferString(`{"name":"tree.xie"}`)
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.True(done)

		assert.Equal(c.StatusCode, http.StatusBadRequest, "error response should be passed by fresh")
		assert.NotNil(c.Body, "error response should be passed by fresh")
		assert.NotNil(c.BodyBuffer, "error response should be passed by fresh")
	})
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
