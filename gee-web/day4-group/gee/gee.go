package gee

import (
	"net/http"
)


// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// Engine should implements the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix string
		middlewares []HandlerFunc
		parent *RouterGroup
		engine *Engine	//all groups share the same engine
	}

	Engine struct{
	router *router
	*RouterGroup	//匿名嵌套结构体，因此Engine能通过.方法直接调用RouterGroup的方法
	groups []*RouterGroup	//all group
	}
)


// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router : newRouter()}
	engine.RouterGroup = &RouterGroup{engine : engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix : group.prefix + prefix,
		parent: group,
		engine : engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern,handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := newContext(writer, request)
	engine.router.handle(context)
}



