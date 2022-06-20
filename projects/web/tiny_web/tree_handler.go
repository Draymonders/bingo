package main

import (
	"fmt"
	"net/http"
	"strings"
)

type httpMethod string

type TreeHandler struct {
	methodRoots map[httpMethod]*node
}

func NewHandlerBasedOnTree() *TreeHandler {
	return &TreeHandler{methodRoots: make(map[httpMethod]*node)}
}

type node struct {
	Path     string     // 当前节点的path
	handleFn HandleFunc // 等价于val
	children []*node
}

func newNode(path string, handleFn HandleFunc) *node {
	return &node{
		Path:     path,
		children: make([]*node, 0, 2),
		handleFn: handleFn,
	}
}

func (s *TreeHandler) Route(method string, path string, handleFunc HandleFunc) {
	var root *node
	root = s.methodRoots[httpMethod(method)]
	if root == nil {
		root = newNode("/", nil)
		s.methodRoots[httpMethod(method)] = root
	}

	root.buildTreeRoute(path, handleFunc)
}

func (s *TreeHandler) ServeHTTP(ctx *Context) {
	if ctx == nil {
		return
	}
	methodRoot := s.methodRoots[httpMethod(ctx.R.Method)]
	if methodRoot == nil {
		ctx.W.WriteHeader(http.StatusBadRequest)
		return
	}
	n := methodRoot.findTreeNode(ctx.R.URL.Path)
	fmt.Printf("treeNode %+v\n", n)
	if n == nil || n.handleFn == nil {
		ctx.W.WriteHeader(http.StatusBadRequest)
		return
	}
	n.handleFn(ctx)
}

func (n *node) findTreeNode(path string) *node {
	paths := rewritePath(path)
	cur := n
	for _, pth := range paths {
		var subNode *node
		for _, son := range cur.children {
			if son.Path == pth {
				subNode = son
				break
			}
		}
		if subNode == nil {
			return subNode
		}
		cur = subNode
	}
	return cur
}

func (n *node) buildTreeRoute(path string, handleFn HandleFunc) {
	// /user
	paths := rewritePath(path)
	cur := n
	for _, pth := range paths {
		var subNode *node
		// 遍历每一个sub，查看是否相等
		for _, son := range cur.children {
			if son.Path == pth {
				subNode = son
				break
			}
		}
		if subNode == nil {
			subNode = newNode(pth, nil)
			n.children = append(n.children, subNode)
		}
		cur = subNode
	}
	cur.handleFn = handleFn
}

func rewritePath(path string) []string {
	paths := strings.Split(path, "/")
	if len(paths) > 0 && paths[0] == "" {
		paths = paths[1:]
	}
	return paths
}
