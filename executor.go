package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type Executor struct {
	tree *UrlsTree
}

func getHeadAndTail(urlPath string) (string, string) {
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

func FindHandlerAndHandleHelper(currentNode *ResourceNode, currentReqPath string, params map[string]string) { //this is the logic for execute; insertHelper should accept type urlPath
	head, _ := getHeadAndTail(currentReqPath)
	if head == "" {
		//handle the request using the handler inside of the current node
		handlerBuilder := currentNode.resource.HandlerBuilder()
		handlerBuilder(params)
	}

	currentPathEntry := currentNode.resource
	switch currentPathEntry.(type) {
	case *ResourceIdentifier:
		//save the entry in the params somehow; prolly using the name of the path entry for now
	case *ResourceCollection:
		//start move on
	case *ResourceSentinel:
		//maybe not needed at all
	default:
		log.Fatal("not possible")
	}
}

func (e *Executor) NewExecutor(tree *UrlsTree) {
	e.tree = tree
}

func (e *Executor) FindHandlerAndHandle(pathInRequest string, w http.ResponseWriter) {
	sentinelNode := e.tree.root
	FindHandlerAndHandleHelper(sentinelNode, pathInRequest, map[string]string{})
}

func (e *Executor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestUrlPath := req.URL.Path
	e.FindHandlerAndHandle(requestUrlPath, w)
}
