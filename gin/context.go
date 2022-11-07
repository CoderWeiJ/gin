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
	Params map[string]string // 提供对路由参数的访问
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

// newContext 创建Context实例
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		// index 记录当前执行到第几个中间件，调用Next()方法将控制权交给下一个中间件
		// 然后再从后往前，调用每个中间件在 Next 方法之后定义的部分。
		index: -1,
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

// Param 获取存储的路由参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Next 中间件放行，按顺序执行所有的中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
