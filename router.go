package router

import (
	"net/http"
	"path"

	"github.com/asmir-a/gorestrouter/endpointparser"
	"github.com/asmir-a/gorestrouter/executor"
	"github.com/asmir-a/gorestrouter/resource"
	"github.com/asmir-a/gorestrouter/tree"
)

type Router struct {
	endpoints []resource.Url
}

func (router *Router) Handle(endpointPath string, handlerBuilder func(map[string]string) http.Handler) {
	url := endpointparser.ParseToUrl(endpointPath, handlerBuilder)
	router.endpoints = append(router.endpoints, url)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	urlsTree := tree.NewUrlsTree(router.endpoints) //the abstraction boundaries are bit weird right now. First, router should contain an executor inside of it. Also, the router should be able to insert a path to executor after every handler is called. Router prolly should not have access to tree variable. The serveHttp prolly just should call FindHandlerAndHandle
	requestExecutor := executor.NewExecutor(urlsTree)

	endpointPath := req.URL.Path
	endpointPath = path.Clean(endpointPath)

	requestExecutor.FindHandlerAndHandle(w, req, endpointPath)
}
