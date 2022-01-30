package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
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
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
	}
)


// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router : newRouter()}
	engine.RouterGroup = &RouterGroup{engine : engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
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

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	context := newContext(writer, request)
	context.handlers = middlewares	//
	context.engine = engine	//
	engine.router.handle(context)
}


//createStaticHandler create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absoluteReqPath := path.Join(group.prefix, relativePath)
	fileServer := http.FileServer(fs)
	fileServer = http.StripPrefix(absoluteReqPath, fileServer)
	return func(c *Context) {
		filepath := c.Param("filepath")
		if _, err := fs.Open(filepath); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//Static serve static files
//r.Static("/assets", "./static")
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "*filepath")
	group.GET(urlPattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates =template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

