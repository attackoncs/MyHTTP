package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc //支持的middlerware
		parent      *RouterGroup  //支持嵌套
		engine      *Engine       //所有组共享一个engine实例
	}
	//Engine is handler for all requests
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
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

func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	//查看http.ListenAndServe中的第二个参数handler，它是interface，只要实现它的方法ServeHTTP的接口，都能被强制转换为接口类型，因此
	//下面等价于log.Fatal(http.ListenAndServe(":9999", (http.Handler)(engine)))，然后ListenAndServe方法里面会去调用 handler.ServeHTTP() 方法
	return http.ListenAndServe(addr, engine)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//该函数定义在http.ListenAndServe的第二个参数handler interface中，参数分别是ResponseWriter和Request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
