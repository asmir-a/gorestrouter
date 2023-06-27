package resource

import "net/http"

type Resource interface {
	Name() string
	HandlerBuilder() func(map[string]string) http.Handler
}

type Url []Resource
