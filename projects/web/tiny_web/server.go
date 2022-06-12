package main

import (
	"net/http"
)

type IRoute interface {
	Route(method string, path string, handleFunc HandleFunc)
}

type IServer interface {
	IRoute

	Start(address string) error // 服务启动
}

// 设计技巧: struct 是否实现对应 interface
var _ IServer = &Server{}

type Server struct {
	handler Handler
	chains  Filter
}

func NewServerWithHandler(handler Handler, builders ...FilterBuilder) *Server {
	wrapF := handler.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		builder := builders[i]
		wrapF = builder(wrapF)
	}

	return &Server{
		handler: handler,
		chains:  wrapF,
	}
}

func (s *Server) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		s.chains(ctx)
	})

	return http.ListenAndServe(address, nil)
}

func (s *Server) Route(method string, path string, handleFunc HandleFunc) {
	s.handler.Route(method, path, handleFunc)
}
