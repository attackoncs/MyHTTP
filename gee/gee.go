package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

//Engine is handler for all requests
type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router[method+"-"+pattern] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	//查看http.ListenAndServe中的第二个参数handler，它是interface，只要实现它的方法ServeHTTP的接口，都能被强制转换为接口类型，因此
	//下面等价于log.Fatal(http.ListenAndServe(":9999", (http.Handler)(engine)))，然后ListenAndServe方法里面会去调用 handler.ServeHTTP() 方法
	return http.ListenAndServe(addr, engine)
}

//该函数定义在http.ListenAndServe的第二个参数handler interface中，参数分别是ResponseWriter和Request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 not found:%s\n", req.URL.Path)
	}
}
