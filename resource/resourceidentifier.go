package resource

import "net/http"

type ResourceIdentifier struct {
	identifierName string
	handlerBuilder func(map[string]string) http.Handler //in langlearn project, this should be httperrors.HandlerWithHttpError
}

func (pe *ResourceIdentifier) Name() string {
	return pe.identifierName
}

func (pe *ResourceIdentifier) HandlerBuilder() func(map[string]string) http.Handler {
	return pe.handlerBuilder
}

func NewResourceIdentifier(identifierName string, handlerBuilder func(params map[string]string) http.Handler) *ResourceIdentifier {
	return &ResourceIdentifier{
		identifierName: identifierName,
		handlerBuilder: handlerBuilder,
	}
}
