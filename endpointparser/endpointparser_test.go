package endpointparser

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/asmir-a/gorestrouter/resource"
)

func TestAttacherGeneral(t *testing.T) {
	endpointPathTest := "/[username]/wordgame/stats"
	handlerBuilderTest := func(params map[string]string) http.Handler {
		username := params["username"]
		handler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handler)
	}
	urlTest := ParseToUrl(endpointPathTest, handlerBuilderTest)

	fmt.Println(urlTest)

	resourcesCount := len(urlTest)
	if resourcesCount != 3 {
		t.Fatal("the count is supposed to be 3, but is: ", resourcesCount)
	}

	firstResourceName := urlTest[0].Name()
	if firstResourceName != "username" {
		t.Fatal("the name of the first resource is supposed to be username, but it is: ", firstResourceName)
	}
	switch urlTest[0].(type) {
	case *resource.ResourceIdentifier:
	default:
		t.Fatal("the type of the first resource is supposed to be idenitifer")
	}

	secondResourceName := urlTest[1].Name()
	if secondResourceName != "wordgame" {
		t.Fatal("the name of the second resource is supposed to be wordgame, but it is: ", secondResourceName)
	}
	switch urlTest[1].(type) {
	case *resource.ResourceCollection:
	default:
		t.Fatal("the type of the resource is supposed to be collection")

	}

	thirdResourceName := urlTest[2].Name()
	if thirdResourceName != "stats" {
		t.Fatal("the name of the third resource is supposed to be stats, but it is: ", thirdResourceName)
	}
	switch urlTest[2].(type) {
	case *resource.ResourceCollection:
	default:
		t.Fatal("the type of the resource is supposed to be collection")
	}
}
