package main

import "net/http"

type ResourceIdentifier struct {
	name           string
	handlerBuilder func(map[string]string) http.Handler //in langlearn project, this should be httperrors.HandlerWithHttpError
}

func (pe *ResourceIdentifier) Name() string {
	return pe.name
}

func (pe *ResourceIdentifier) HandlerBuilder() func(map[string]string) http.Handler {
	return pe.handlerBuilder
}
