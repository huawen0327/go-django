package golee

import (
	"net/http"
)

// HandlerFunc 将请求定义为一个结构体
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine ServeHTTP的结构体
type Engine struct {
	router map[string]HandlerFunc
}

// New 设置路径
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// Path 添加路径
func (engine *Engine) Path(pattern string, handler HandlerFunc) {
	engine.router[pattern] = handler
}
