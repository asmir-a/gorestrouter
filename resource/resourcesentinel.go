package resource

import "net/http"

type ResourceSentinel struct {
	name           string
	handlerBuilder func(map[string]string) http.Handler
}

func (pe *ResourceSentinel) Name() string {
	return pe.name
}

func (pe *ResourceSentinel) HandlerBuilder() func(map[string]string) http.Handler {
	return pe.handlerBuilder
}

func NewResourceSentinel() *ResourceSentinel {
	return &ResourceSentinel{
		name:           "",
		handlerBuilder: nil,
	}
}
