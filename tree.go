package router

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	SEPARATOR string = "/"
)

type node struct {
	part     string
	handlers map[string]http.Handler
	children map[string]*node
}

func newNode(part string) *node {
	return &node{
		part:     part,
		handlers: make(map[string]http.Handler),
		children: make(map[string]*node),
	}
}

type tree struct {
	root *node
}

func NewTree() *tree {
	return &tree{root: newNode(SEPARATOR)}
}

func (t *tree) i(method, pattern string, handler http.Handler) {
	curNode := t.root
	if pattern != SEPARATOR {
		for _, part := range split(pattern) {
			next, ok := curNode.children[part]
			if !ok {
				next = newNode(part)
				curNode.children[part] = next
			}
			curNode = next
		}
	}
	if _, ok := curNode.handlers[method]; ok {
		panic(fmt.Sprintf("%s: handler already registered on [%s %s]", MOD, method, pattern))
	}
	curNode.handlers[method] = handler
}

func (t *tree) lookup(method, path string) http.Handler {
	curNode := t.root
	if path != SEPARATOR {
		for _, part := range split(path) {
			next, ok := curNode.children[part]
			if !ok {
				return err404
			}
			curNode = next
		}
	}

	hLen := len(curNode.handlers)
	if hLen == 0 {
		return err404
	}
	if h, ok := curNode.handlers[method]; ok {
		return h
	}
	methods := make([]string, 0, hLen)
	for m := range curNode.handlers {
		methods = append(methods, m)
	}
	return err405{allowedMethods: methods}
}

func split(pat string) []string {
	parts := strings.Split(pat, SEPARATOR)
	var s []string
	for _, p := range parts {
		if p != "" {
			s = append(s, p)
		}
	}
	return s
}
