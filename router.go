package golee

import (
	"errors"
	"net/http"
	"strings"
)

// HandlerFunc 将请求定义为一个结构体
type HandlerFunc func(*Context)

// node节点
type node struct {
	// pattern     string      // 待匹配路由，例如 /p/:lang
	pathNode    string        // 路由中的一部分，例如 :lang
	children    []*node       // 子节点，例如 [doc, tutorial, intro]
	isWild      bool          // 是否精确匹配，part开头为 * 时为true
	handlerFunc []HandlerFunc // 配置函数
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
func (n *node) insert(pathArr []string, handler []HandlerFunc, level int) {
	// 判断是否是 * 开头的路径
	if pathArr[level][0] == '*' {
		child := &node{pathNode: pathArr[level], handlerFunc: handler, isWild: true}
		n.children = append(n.children, child)
		return
	}
	// 判断是否为最后一个
	if level+1 == len(pathArr) {
		child := &node{pathNode: pathArr[level], handlerFunc: handler}
		n.children = append(n.children, child)
		return
	}
	var child *node
	// 寻找当前路径是否已经存在
	for _, c := range n.children {
		if c.pathNode == pathArr[level] {
			child = c
			break
			// fmt.Println(child)
		}
	}
	if child == nil {
		child = &node{pathNode: pathArr[level]}
	}
	// child := &node{pathNode: pathArr[level]}
	n.children = append(n.children, child)
	child.insert(pathArr, handler, level+1)
}

// 添加首节点
func (n *node) addNode(pathArr []string, handler []HandlerFunc) {
	n.pathNode = "/"
	if pathArr[1] == "" {
		n.handlerFunc = handler
		return
	}
	n.insert(pathArr, handler, 1)
	return
}

// Path 添加路径
func (router *Router) Path(pattern string, handler ...HandlerFunc) {
	pathArr := strings.Split(pattern, "/")
	pathArr[0] = "/"
	var handlerArr []HandlerFunc
	for _, h := range handler {
		handlerArr = append(handlerArr, h)
	}
	router.nodes.addNode(pathArr, handlerArr)
}

// 查找最终Node的HandlerFunc
func (n *node) matchNode(c *Context, level int) ([]HandlerFunc, error) {
	pathArr := strings.Split(c.Path, "/")
	if len(pathArr) == level+1 {
		for _, child := range n.children {
			if (pathArr[level] == child.pathNode || child.pathNode[0] == ':' || child.pathNode[0] == '*') && child.handlerFunc != nil {
				c.Parm = pathArr[level]
				handler := child.handlerFunc
				return handler, nil
			}
		}
		return nil, errors.New("Not Found Path")
	}
	for _, child := range n.children {
		if child.pathNode[0] == '*' {
			c.Parm = pathArr[level]
			handler := child.handlerFunc
			return handler, nil
		}
		if child.pathNode[0] == ':' {
			c.Parm = pathArr[level]
			handler, err := child.matchNode(c, level+1)
			return handler, err
		}
		if pathArr[level] == child.pathNode {
			handler, err := child.matchNode(c, level+1)
			return handler, err
		}
	}
	return nil, errors.New("Not Found Path")
}

// 查找路径
func (router *Router) match(c *Context) ([]HandlerFunc, error) {
	pathArr := strings.Split(c.Path, "/")
	if len(pathArr) == 2 && pathArr[1] == "" {
		handler := router.nodes.handlerFunc
		if handler == nil {
			return nil, errors.New("Not Found Path")
		}
		return handler, nil
	}
	handler, err := router.nodes.matchNode(c, 1)
	return handler, err
}

func (router *Router) handle(c *Context) {
	// key := c.Path
	handler, err := router.match(c)
	if err != nil {
		c.StringResponse(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	} else {
		c.handler = handler
		c.run()
	}
}
