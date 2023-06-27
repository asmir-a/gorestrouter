package restrouter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestRouterGeneral(t *testing.T) {
	router := &Router{}
	handlerBuilder := func(params map[string]string) http.Handler {
		username := params["username"]
		handler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handler)
	}
	router.Handle("/[username]/wordgame/stats", handlerBuilder)

	server := httptest.NewServer(router)
	baseServerUrl := server.URL
	defer server.Close()

	client := http.Client{}
	response, err := client.Get(baseServerUrl + "/asmir/wordgame/stats")
	if err != nil {
		t.Fatal("request could not be sent to server: ", err)
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("could not read response body: ", err)
	}
	bodyString := string(bodyBytes)
	if bodyString != "asmir" {
		t.Fatal("the server should have sent asmir to the client, but sent", bodyString)
	}
}
