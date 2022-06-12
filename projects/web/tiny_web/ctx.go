package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request

	bodyBytes []byte
	mux       sync.RWMutex
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}

func (c *Context) ReadJson(data interface{}) (err error) {
	if err = c.ReadBodyBytes(); err != nil {
		return err
	}
	return json.Unmarshal(c.bodyBytes, data)
}

func (c *Context) ReadBodyBytes() (err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if len(c.bodyBytes) > 0 {
		return
	}
	c.bodyBytes, err = ioutil.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) String() string {
	req := fmt.Sprintf("req: [path=%v, method=%v, body=%v]", c.R.URL.Path, c.R.Method, string(c.bodyBytes))
	resp := fmt.Sprintf("resp: [header=%v]", c.W.Header())

	return req + " " + resp

}
