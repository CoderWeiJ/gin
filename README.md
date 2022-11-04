# HTTP 基础
- 定义一个`HandlerFunc`类型
```go
type HandlerFunc func(w http.ResponseWriter, r *http.Request)
```
- 定义一个Engine结构体，实现路由注册表功能、gin 的`Run`和`ServeHttp`方法
```go
// addRoute 添加路由
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

// GET 注册到路由映射表
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST 注册到路由映射表
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Run 定义一个开启http服务的方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 解析请求的路径，查找路由映射表
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

```

# Context 上下文
> 封装`*http.Request`和`http.Responser`的方法，扩展额外功能，比如：存放动态路由，注册中间件，保存上下文等
 
gee/context.go
```go
package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 起别名
type H map[string]interface{}

// Context 包含最基本的 http.ResponseWriter 和 *http.Request
// 提供对 Method 和 Path 的直接访问
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

// newContext 创建Context实例
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm 获取表单信息
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取url后携带的参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置请求状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置头部信息
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回字符串信息
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 返回JSON数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 直接返回数据和状态码
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 把HTML文本返回
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
```

# 路由（Router）
将和路由相关的方法和结构体提取出来，方便扩展新功能

# 前缀树路由
树节点应该存储的信息
```go
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}
```