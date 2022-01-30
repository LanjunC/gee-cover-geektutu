package gee

import (
	"log"
	"net/http"
	"strings"
)


type router struct {
	roots map[string]*node	//GET、POST等请求方式的根节点
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots : make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	split := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _,item := range split {
		if item != "" {	//go的split会产生""
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, pattern)

	_,ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root,ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
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
		c.Params = params	//注意这里需要将动态路由的参数给解析了，也可以在newContext时解析
		key := c.Method + "-" + n.pattern	//注意这里的key
		c.handlers = append(c.handlers, r.handlers[key])	//直到这时才把本来的handler加入中间件调用链
	}else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	c.Next()
}