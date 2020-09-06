package golee

import "strings"

// Group 分组管理
type Group struct {
	prefix string
	// router       *Router
	HandlerFuncs []HandlerFunc
	engine       *Engine
}

// Path 集成自router的函数
func (g *Group) Path(pattern string, handler ...HandlerFunc) {
	pattern = g.prefix + pattern
	pathArr := strings.Split(pattern, "/")
	pathArr[0] = "/"
	var handlerArr []HandlerFunc
	handlerArr = append(handlerArr, g.HandlerFuncs...)
	for _, h := range handler {
		handlerArr = append(handlerArr, h)
	}
	g.engine.Router.nodes.addNode(pathArr, handlerArr)
}

// Use 中间件使用函数
func (g *Group) Use(handler ...HandlerFunc) {
	for _, h := range handler {
		g.HandlerFuncs = append(g.HandlerFuncs, h)
	}
}
