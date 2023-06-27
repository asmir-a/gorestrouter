package endpointparser

import (
	"net/http"
	"path"
	"strings"

	"github.com/asmir-a/gorestrouter/resource"
)

func isIdentifierName(token string) bool {
	if len(token) == 0 { //some validation functionality is required
		return false
	}
	if token[0] == '[' && token[len(token)-1] == ']' {
		return true
	}
	return false
}

func extractIdentifierName(token string) string {
	return token[1 : len(token)-1]
}

func createResource(resourceName string, handlerBuilder func(map[string]string) http.Handler) resource.Resource {
	if resourceName == "" { //maybe should have a method to check if the name is valid
		return nil
	}
	if isIdentifierName(resourceName) {
		identifierName := extractIdentifierName(resourceName)
		return resource.NewResourceIdentifier(identifierName, handlerBuilder)
	}
	return resource.NewResourceCollection(resourceName, handlerBuilder)
}

func ParseToUrl(endpointUrl string, handlerBuilder func(map[string]string) http.Handler) resource.Url {
	url := resource.Url{}
	endpointUrl = path.Clean(endpointUrl)
	resourceNames := strings.Split(endpointUrl, "/")
	for resourceIndex, resourceName := range resourceNames {
		if resourceIndex == len(resourceNames)-1 {
			newResource := createResource(resourceName, handlerBuilder)
			if newResource != nil {
				url = append(url, newResource)
			}
			continue
		}
		newResource := createResource(resourceName, nil)
		if newResource != nil {
			url = append(url, createResource(resourceName, nil))
		}
	}
	return url
}
