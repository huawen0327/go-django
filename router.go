package golee

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc 将请求定义为一个结构体
type HandlerFunc func(*Context)

// node节点
type node struct {
	// pattern     string      // 待匹配路由，例如 /p/:lang
	pathNode    string      // 路由中的一部分，例如 :lang
	children    []*node     // 子节点，例如 [doc, tutorial, intro]
	isWild      bool        // 是否精确匹配，part 含有 : 或 * 时为true
	handlerFunc HandlerFunc // 配置函数
}

// Router ServeHTTP的结构体
type Router struct {
	paths map[string]HandlerFunc
	nodes *node
}

// New 设置路径
func newRouter() *Router {
	return &Router{
		paths: make(map[string]HandlerFunc),
		nodes: new(node),
	}
}

// 添加节点
func (n *node) insert(pathArr []string, handler HandlerFunc) {
	fmt.Println("point")
}

// 添加首节点
func (n *node) addNode(pathArr []string, handler HandlerFunc) {
	n.pathNode = "/"
	if pathArr[1] == "" {
		n.handlerFunc = handler
		return
	}
	n.insert(pathArr, handler)
	return
}

// Path 添加路径
func (router *Router) Path(pattern string, handler HandlerFunc) {
	pathArr := strings.Split(pattern, "/")
	pathArr[0] = "/"
	router.nodes.addNode(pathArr, handler)
	// for index, pathNode := range pathArr {
	// 	insert(pathArr, pathNode, index)
	// }
	// router.paths[pattern] = handler
	// reqNode := &node{
	// 	part:   "b",
	// 	isWild: true,
	// }
	// router.nodes = append(router.nodes, reqNode)
}

func (router *Router) handle(c *Context) {
	key := c.Path
	if handler, ok := router.paths[key]; ok {
		handler(c)
	} else {
		c.StringResponse(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
