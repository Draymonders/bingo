package main

import (
	"fmt"
	"net/http"
	"sync"
)

type HandleFunc func(ctx *Context)

type Handler interface {
	IRoute

	ServeHTTP(*Context)
}

// MapHandler 使用 sync.Map 来存储 route 信息
type MapHandler struct {
	Routers sync.Map
}

func NewHandlerBaseOnMap() Handler {
	return &MapHandler{}
}

func (s *MapHandler) Route(method string, path string, handleFunc HandleFunc) {
	s.Routers.Store(s.reqKey(method, path), handleFunc)
}

func (s *MapHandler) ServeHTTP(ctx *Context) {
	if ctx == nil {
		return
	}
	method, path := ctx.R.Method, ctx.R.URL.Path
	key := s.reqKey(method, path)
	val, ok := s.Routers.Load(key)
	if !ok {
		ctx.W.WriteHeader(http.StatusBadRequest)
		return
	}
	val.(HandleFunc)(ctx)
}

func (s *MapHandler) reqKey(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}
