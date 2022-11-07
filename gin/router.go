package gin

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // 存储每种请求方式的Trie树根节点
	handlers map[string]HandlerFunc // 存储每种请求方式的HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析请求的路径，分解成一个一个节点
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0) // 每一层的节点
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 遇到动态路由，退出
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{} // 找不到，则创建根节点
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler

	log.Printf("Router %4s - %s", method, pattern)
}

// getRoute 解析':'和'*'两种匹配符的参数
// 返回一个map
// eg:
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // 查找根节点

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		// 在调用匹配到的handler之前，将解析出来的路由参数赋值给c.Params，就能在handlers里面访问到具体的值
		c.Params = params
		// 先执行中间件
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
