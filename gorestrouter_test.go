package gorestrouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicFunctionality(t *testing.T) {
	router := &Router{}
	wordsHandlerBuilder := func(params map[string]string) http.Handler {
		username := params["username"]
		handler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handler)
	}
	router.Handle("/[username]/wordgame/words", wordsHandlerBuilder)
	statsHandlerBuilder := func(params map[string]string) http.Handler {
		username := params["username"]
		handler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handler)
	}
	router.Handle("/wordgame/[username]/stats", statsHandlerBuilder)

	requestToWords := httptest.NewRequest("GET", "/asmir/wordgame/words", nil)
	requestToStats := httptest.NewRequest("GET", "/wordgame/nina/stats", nil)

	responseFromWords := httptest.NewRecorder()
	responseFromStats := httptest.NewRecorder()

	router.ServeHTTP(responseFromWords, requestToWords)
	router.ServeHTTP(responseFromStats, requestToStats)

	responseFromWordsString := responseFromWords.Body.String()
	responseFromStatsString := responseFromStats.Body.String()
	if responseFromWordsString != "asmir" {
		t.Fatal("wanted: ", "asmir", ", but got: ", responseFromWordsString)
	}
	if responseFromStatsString != "nina" {
		t.Fatal("wanted: ", "nina", ", but got: ", responseFromStatsString)
	}
}
