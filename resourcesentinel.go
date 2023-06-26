package main

import "net/http"

type ResourceSentinel struct {
	name         string
	setUpHandler func(map[string]string) http.Handler
}

func (pe *ResourceSentinel) Name() string {
	return pe.name
}

func (pe *ResourceSentinel) HandlerBuilder() func(map[string]string) http.Handler {
	return pe.setUpHandler
}
