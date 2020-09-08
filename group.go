package golee

import "strings"

// Group 分组管理
type Group struct {
	prefix string
	// router       *Router
	handler []HandlerFunc
	engine  *Engine
}

// Path 集成自router的函数
func (g *Group) Path(pattern string, handler ...HandlerFunc) {
	if g.prefix != "/" {
		pattern = g.prefix + pattern
	}
	if pattern == "/" {
		pattern = g.prefix
	}
	pathArr := strings.Split(pattern, "/")
	pathArr[0] = "/"
	var handlerArr []HandlerFunc
	handlerArr = append(handlerArr, g.handler...)
	handlerArr = append(handlerArr, handler...)
	// for _, h := range handler {
	// 	handlerArr = append(handlerArr, h)
	// }
	g.engine.Router.nodes.addNode(pathArr, handlerArr)
}

// SetGroup 中间件使用函数
func (g *Group) SetGroup(prefix string, handler ...HandlerFunc) {
	g.prefix = prefix
	// for _, h := range handler {
	g.handler = append(g.handler, handler...)
	// }
}
