package executor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asmir-a/gorestrouter/resource"
	"github.com/asmir-a/gorestrouter/tree"
)

func TestGetHeadAndTail(t *testing.T) {
	head1, tail1 := getHeadAndTail("")
	if head1 != "" && tail1 != "" {
		t.Fatalf("head1 and tail1 are supposed to be empty, but are: %s and %s", head1, tail1)
	}
	head2, tail2 := getHeadAndTail("/")
	if head2 != "" && tail2 != "" {
		t.Fatalf("head2 and tail2 are supposed to be empty, but are: %s and %s", head2, tail2)
	}
	head3, tail3 := getHeadAndTail("/whatever")
	if head3 != "whatever" && tail3 != "" {
		t.Fatalf("head3 and tail3 are supposed to be whatever and empty respectively, but are %s and %s", head3, tail3)
	}
	head4, tail4 := getHeadAndTail("/whatever/")
	if head4 != "whatever" && tail4 != "" {
		t.Fatalf("head4 and tail4 are supposed to be whatever and empty respectively, but are %s and %s", head4, tail4)
	}
	head5, tail5 := getHeadAndTail("/whatever/something")
	if head5 != "whatever" && tail5 != "something" {
		t.Fatalf("head5 and tail5 are supposed to be whatever and something, but are %s and %s", head5, tail5)
	}
}

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
