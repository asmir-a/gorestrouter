package resource

import "net/http"

type ResourceCollection struct {
	collectionName string
	handlerBuilder func(map[string]string) http.Handler
}

func (pe *ResourceCollection) Name() string {
	return pe.collectionName
}

func (pe *ResourceCollection) HandlerBuilder() func(map[string]string) http.Handler {
	return pe.handlerBuilder
}

func NewResourceCollection(collectionName string, handlerBuilder func(params map[string]string) http.Handler) *ResourceCollection {
	return &ResourceCollection{
		collectionName: collectionName,
		handlerBuilder: handlerBuilder,
	}
}
