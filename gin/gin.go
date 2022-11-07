package gin

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// HandlerFunc 定义了gin的请求方法
type HandlerFunc func(c *Context)

type (
	RouterGroup struct {
		prefix      string        // 分组前缀
		middlewares []HandlerFunc // 支持中间件
		parent      *RouterGroup  // 支持嵌套
		engine      *Engine
	}

	// Engine 实现了ServeHttp接口
	// 拥有 RouterGroup 所有的能力
	Engine struct {
		*RouterGroup
		router *router // 路由映射表
		groups []*RouterGroup
	}
)

// New gin.Engine构造器
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.Use(Logger())
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use 注册中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	e.router.handle(c)
}

// Logger 全局默认注册的中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
