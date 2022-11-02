package gin

import (
	"fmt"
	"net/http"
)

// HandlerFunc 定义了gin的请求方法
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine 实现了ServeHttp接口
type Engine struct {
	router map[string]HandlerFunc // 路由映射表
}

// New gin.Engine构造器
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

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
