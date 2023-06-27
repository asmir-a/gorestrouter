package main

import "net/http"

type ResourceCollection struct {
	name           string
	handlerBuilder func(map[string]string) http.Handler
}

func (pe *ResourceCollection) Name() string {
	return pe.name
}

func (pe *ResourceCollection) HandlerBuilder() func(map[string]string) http.Handler {
	return pe.handlerBuilder
}
