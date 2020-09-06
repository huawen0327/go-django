package golee

import (
	"net/http"
)

// Engine golee的方法
type Engine struct {
	Router *Router
	Group  *Group
}

// StatusOK OK码
var StatusOK = http.StatusOK

// New engine的构造方法
func New() *Engine {
	engine := &Engine{Router: newRouter()}
	engine.Group = &Group{engine: engine}
	return engine
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP ServeHTTP
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.Router.handle(c)
}
