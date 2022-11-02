package gin

import (
	"net/http"
)

// HandlerFunc 定义了gin的请求方法
type HandlerFunc func(c *Context)

// Engine 实现了ServeHttp接口
type Engine struct {
	router *router // 路由映射表
}

// New gin.Engine构造器
func New() *Engine {
	return &Engine{router: newRouter()}
}

// GET 注册到路由映射表
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

// POST 注册到路由映射表
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

// Run 定义一个开启http服务的方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 解析请求的路径，查找路由映射表
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
