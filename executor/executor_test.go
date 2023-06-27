package executor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asmir-a/gorestrouter/resource"
	"github.com/asmir-a/gorestrouter/tree"
)

func TestExecutorGeneral(t *testing.T) {
	statsHandlerBuilder := func(params map[string]string) http.Handler {
		username := params["username"] //for now this is okay. In the future, the params should be auto passed through somehow
		statsHandler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(statsHandler)
	}
	urlOne := resource.Url{
		resource.NewResourceIdentifier("username", nil),
		resource.NewResourceCollection("wordgame", nil),
		resource.NewResourceCollection("stats", statsHandlerBuilder),
	}
	wordsHandlerBuilder := func(params map[string]string) http.Handler {
		username := params["username"]
		wordsHandler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(wordsHandler)
	}
	urlTwo := resource.Url{
		resource.NewResourceIdentifier("username", nil),
		resource.NewResourceCollection("wordgame", nil),
		resource.NewResourceCollection("words", wordsHandlerBuilder),
	}
	urls := []resource.Url{urlOne, urlTwo}
	urlsTree := tree.NewUrlsTree(urls)
	executor := NewExecutor(urlsTree)

	newRequestOne := httptest.NewRequest("GET", "/asmir/wordgame/stats", nil)
	newRequestTwo := httptest.NewRequest("GET", "/asmir/wordgame/words", nil)

	newRespOne := httptest.NewRecorder()
	newRespTwo := httptest.NewRecorder()

	executor.ServeHTTP(newRespOne, newRequestOne)
	executor.ServeHTTP(newRespTwo, newRequestTwo)

	respBodyOne := newRespOne.Body.String()
	respBodyTwo := newRespTwo.Body.String()

	if respBodyOne != "asmir" {
		t.Fatalf(`the body is %q; wanted body to be %q`, respBodyOne, "asmir")
	}
	if respBodyTwo != "asmir" {
		t.Fatalf(`the body is %q; wanted body to be %q`, respBodyTwo, "asmir")
	}
}
