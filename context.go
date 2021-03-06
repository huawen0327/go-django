package golee

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// H 1
type H map[string]interface{}

// Context 上下文的结构体
type Context struct {
	// origin objects
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	Parm       string
	index      int
	handler    []HandlerFunc
}

// newContext Context的构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  0,
	}
}

// PostForm 获取请求的post字段
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// GET 获取Get请求的字段
func (c *Context) GET(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置Content-Type
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// HTMLResponse Http响应
func (c *Context) HTMLResponse(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// HTMLRender Http响应
func (c *Context) HTMLRender(code int, dict interface{}, pattern ...string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if tmpl, err := template.ParseFiles(pattern...); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	} else {
		tmpl.Execute(c.Writer, dict)
	}
}

// JSONResponse Json响应
func (c *Context) JSONResponse(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// StringResponse String响应
func (c *Context) StringResponse(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Next 下一个函数
func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handler); c.index++ {
		c.handler[c.index](c)
	}
}

func (c *Context) run() {
	for ; c.index < len(c.handler); c.index++ {
		c.handler[c.index](c)
	}
}
