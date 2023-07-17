package executor

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/asmir-a/gorestrouter/resource"
	"github.com/asmir-a/gorestrouter/tree"
)

type Executor struct {
	tree *tree.UrlsTree
}

func getHeadAndTail(urlPath string) (string, string) {
	if urlPath == "/" || urlPath == "" {
		return "", ""
	}
	urlPath = path.Clean("/"+urlPath) + "/"
	if urlPath == "/" { //this assumes that the url is kind of properly formatted; todo: think of what attackers migth try to do
		return "", ""
	}
	slashIndex := strings.Index(urlPath[1:], "/") + 1
	if slashIndex == 0 {
		log.Panic("this is not supposed to happen") //because we are adding a slash to the end of the path
	}
	return urlPath[1:slashIndex], urlPath[slashIndex:]
}

func (e *Executor) FindHandlerAndHandleHelper(
	w http.ResponseWriter,
	req *http.Request,
	currentNode *tree.ResourceNode,
	currentReqUrl string,
	params map[string]string,
) { //this is the logic for execute; insertHelper should accept type urlPath
	head, tail := getHeadAndTail(currentReqUrl)
	if head == "" {
		//handle the request using the handler inside of the current node
		handlerBuilder := currentNode.Resource.HandlerBuilder()
		if handlerBuilder == nil {
			http.Error(w, "resource not found", http.StatusNotFound)
			return
		}
		handler := handlerBuilder(params)
		handler.ServeHTTP(w, req)
		return
	}

	currentResource := currentNode.Resource
	switch currentResource.(type) {
	case *resource.ResourceIdentifier:
		//save the entry in the params somehow; prolly using the name of the path entry for now
		nextNode := currentNode.FindChildWithResourceName(head, params)
		if nextNode == nil {
			http.Error(w, "requested resource is not found", http.StatusNotFound)
			return
		}
		e.FindHandlerAndHandleHelper(w, req, nextNode, tail, params)
	case *resource.ResourceCollection:
		//start move on
		nextNode := currentNode.FindChildWithResourceName(head, params)
		if nextNode == nil {
			http.Error(w, "requested resource is not found", http.StatusNotFound)
			return
		}
		e.FindHandlerAndHandleHelper(w, req, nextNode, tail, params)
	case *resource.ResourceSentinel:
		nextNode := currentNode.FindChildWithResourceName(head, params)
		if nextNode == nil {
			http.Error(w, "request resource is not found", http.StatusNotFound)
			return
		}
		e.FindHandlerAndHandleHelper(w, req, nextNode, tail, params)
	default:
		log.Fatal("not possible")
	}
}

func NewExecutor(tree *tree.UrlsTree) *Executor {
	newExecutor := &Executor{}
	newExecutor.tree = tree
	return newExecutor
}

func (e *Executor) FindHandlerAndHandle(w http.ResponseWriter, req *http.Request, pathInRequest string) {
	sentinelNode := e.tree.Root
	e.FindHandlerAndHandleHelper(w, req, sentinelNode, pathInRequest, map[string]string{})
}

func (e *Executor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestUrlPath := req.URL.Path
	e.FindHandlerAndHandle(w, req, requestUrlPath)
}
