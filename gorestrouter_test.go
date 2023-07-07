package gorestrouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestMultipleParamsAndJson(t *testing.T) {
	type sumResponse struct {
		NumberOne int `json:"numberOne"`
		NumberTwo int `json:"numberTwo"`
		Sum       int `json:"sum"`
	}
	type productResponse struct {
		NumberOne int `json:"numberOne"`
		NumberTwo int `json:"numberTwo"`
		Product   int `json:"product"`
	}
	router := Router{}
	endpointSum := "/numbers/[numberOne]/[numberTwo]/sum"
	endpointProduct := "/numbers/[numberOne]/and/[numberTwo]/product"
	handlerBuilderSum := func(params map[string]string) http.Handler {
		numberOne, err := strconv.Atoi(params["numberOne"])
		if err != nil {
			t.Fatalf("could not convert %q into an int", params["numberOne"])
		}
		numberTwo, err := strconv.Atoi(params["numberTwo"])
		if err != nil {
			t.Fatalf("could not convert %q into an int", params["numberTwo"])
		}
		handler := func(w http.ResponseWriter, req *http.Request) {
			sumData := &sumResponse{
				NumberOne: numberOne,
				NumberTwo: numberTwo,
				Sum:       numberOne + numberTwo,
			}
			sumJsonBytes, err := json.Marshal(sumData)
			if err != nil {
				t.Fatalf("could not encode the struct %q into json", sumData)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(sumJsonBytes)
		}
		return http.HandlerFunc(handler)
	}
	handlerBuilderProduct := func(params map[string]string) http.Handler {
		numberOne, err := strconv.Atoi(params["numberOne"])
		if err != nil {
			t.Fatalf("could not convert %q into an int", params["numberOne"])
		}
		numberTwo, err := strconv.Atoi(params["numberTwo"])
		if err != nil {
			t.Fatalf("could not convert %q into an int", params["numberTwo"])
		}
		handler := func(w http.ResponseWriter, req *http.Request) {
			productData := &productResponse{
				NumberOne: numberOne,
				NumberTwo: numberTwo,
				Product:   numberOne * numberTwo,
			}
			productJsonBytes, err := json.Marshal(productData)
			if err != nil {
				t.Fatalf("could not encode the struct %q into json", productData)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(productJsonBytes)
		}
		return http.HandlerFunc(handler)
	}
	router.Handle(endpointSum, handlerBuilderSum)
	router.Handle(endpointProduct, handlerBuilderProduct)

	requestToSum := httptest.NewRequest("GET", "/numbers/123/321/sum", nil)
	requestToProduct := httptest.NewRequest("GET", "/numbers/12/and/34/product", nil)

	responseFromSum := httptest.NewRecorder()
	responseFromProduct := httptest.NewRecorder()
	router.ServeHTTP(responseFromSum, requestToSum)
	router.ServeHTTP(responseFromProduct, requestToProduct)

	sumJson := responseFromSum.Body.Bytes()
	productJson := responseFromProduct.Body.Bytes()

	var sumData sumResponse
	json.Unmarshal(sumJson, &sumData)
	if sumData.NumberOne != 123 {
		t.Fatalf("expected the first number to be 123 but it is %q", sumData.NumberOne)
	} else if sumData.NumberTwo != 321 {
		t.Fatalf("expected the second number to be 321 but it is %q", sumData.NumberTwo)
	} else if sumData.Sum != 123+321 {
		t.Fatalf("expected the sum to be %q but it is %q", 123+321, sumData.Sum)
	}

	var productData productResponse
	json.Unmarshal(productJson, &productData)
	if productData.NumberOne != 12 {
		t.Fatalf("expected the first number to be 123 but it is %q", productData.NumberOne)
	} else if productData.NumberTwo != 34 {
		t.Fatalf("expected the second number to be 321 but it is %q", productData.NumberTwo)
	} else if productData.Product != 12*34 {
		t.Fatalf("expected the product to be %q but it is %q", 12*34, productData.Product)
	}
}

func TestNotFound(t *testing.T) {
	router := Router{}
	endpointWords := "/[username]/wordgame/words"
	handlerBuilderWords := func(params map[string]string) http.Handler {
		username := params["username"]
		handler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handler)
	}
	router.Handle(endpointWords, handlerBuilderWords)

	requestToWords := httptest.NewRequest("GET", "/asmir/sentencegame/words", nil)
	responseFromWords := httptest.NewRecorder()

	router.ServeHTTP(responseFromWords, requestToWords)

	if responseFromWords.Result().StatusCode != http.StatusNotFound {
		t.Fatal("expected a response with not found status")
	}
}

func TestAppendRouter(t *testing.T) {
	endpointLogin := "/api/auth/login"
	handlerLogin := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("login"))
	})

	mux := http.ServeMux{}
	mux.Handle(endpointLogin, handlerLogin)

	endpointUsername := "/api/users/[username]"
	handlerBuilderUsername := func(params map[string]string) http.Handler {
		username := params["username"]
		handlerUsername := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handlerUsername)
	}
	usernameRouter := &Router{}
	usernameRouter.Handle(endpointUsername, handlerBuilderUsername)
	mux.Handle("/api/users/", usernameRouter)

	requestToLogin := httptest.NewRequest("GET", "/api/auth/login", nil)
	responseFromLogin := httptest.NewRecorder()
	mux.ServeHTTP(responseFromLogin, requestToLogin)
	responseFromLoginString := responseFromLogin.Body.String()
	if responseFromLoginString != "login" {
		t.Fatalf("wanted to receive \"login\", but received: %q", responseFromLoginString)
	}

	requestToUsername := httptest.NewRequest("GET", "/api/users/asmir", nil)
	responseFromUsername := httptest.NewRecorder()
	mux.ServeHTTP(responseFromUsername, requestToUsername)
	responseFromUsernameString := responseFromUsername.Body.String()
	if responseFromUsernameString != "asmir" {
		t.Fatalf("wanted to receive \"asmir\", but received: %q", responseFromUsernameString)
	}
}

func TestWithStripPrefix(t *testing.T) {
	endpointLogin := "/api/auth/login"
	handlerLogin := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("login"))
	})
	mux := http.ServeMux{}
	mux.Handle(endpointLogin, handlerLogin)

	endpointUsername := "/wordgame/[username]"
	handlerBuilderUsername := func(params map[string]string) http.Handler {
		username := params["username"]
		handlerUsername := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(username))
		}
		return http.HandlerFunc(handlerUsername)
	}

	router := &Router{}
	router.Handle(endpointUsername, handlerBuilderUsername)
	mux.Handle("/api/users/", http.StripPrefix("/api/users/", router)) //todo: should try to implement the ability to pass the router wihtout the stripprefix

	requestToLogin := httptest.NewRequest("GET", "/api/auth/login", nil) //should have a test that will test the correctness for multiple paths. And then not to put multiple paths when they are not need like in here
	responseToLogin := httptest.NewRecorder()
	mux.ServeHTTP(responseToLogin, requestToLogin)
	responseToLoginString := responseToLogin.Body.String()
	if responseToLoginString != "login" {
		t.Fatalf("wanted to receive \"login\", but received: %q", responseToLoginString)
	}

	requestToUsername := httptest.NewRequest("GET", "/api/users/wordgame/asmir", nil)
	responseToUsername := httptest.NewRecorder()
	mux.ServeHTTP(responseToUsername, requestToUsername)
	responseToUsernameString := responseToUsername.Body.String()
	if responseToUsernameString != "asmir" {
		t.Fatalf("wanted to receive \"asmir\", but received: %q", responseToUsernameString)
	}
}

func TestNotFoundTwo(t *testing.T) {
	mux := http.ServeMux{}

	router := &Router{}
	handlerBuilderSubmit := func(params map[string]string) http.Handler {
		username := params["username"]
		handlerSubmit := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(fmt.Sprintf("submit for username: %s", username)))
		}
		return http.HandlerFunc(handlerSubmit)
	}
	router.Handle("/entries/users/[username]/submit", handlerBuilderSubmit)

	mux.Handle("/api/wordgame/", http.StripPrefix("/api/wordgame/", router))

	requestToSubmit := httptest.NewRequest("GET", "/api/wordgame/entries/users/koala/submit", nil)
	responseFromSubmit := httptest.NewRecorder()
	mux.ServeHTTP(responseFromSubmit, requestToSubmit)

	want := "submit for username: koala"
	got := responseFromSubmit.Body.String()
	if got != want {
		t.Fatalf("expected %s but got %s", want, got)
	}

	requestToWordgame := httptest.NewRequest("GET", "/api/wordgame/", nil)
	responseFromWordgame := httptest.NewRecorder()
	mux.ServeHTTP(responseFromWordgame, requestToWordgame)

	wantedCode := http.StatusNotFound
	gotCode := responseFromWordgame.Result().StatusCode
	if gotCode != wantedCode {
		t.Fatalf("wanted code %d but got code %d", wantedCode, gotCode)
	}

	requestToEntries := httptest.NewRequest("GET", "/api/wordgame/", nil)
	responseFromEntries := httptest.NewRecorder()
	mux.ServeHTTP(responseFromEntries, requestToEntries)

	wantedCode = http.StatusNotFound
	gotCode = responseFromWordgame.Result().StatusCode
	if gotCode != wantedCode {
		t.Fatalf("wanted code %d but got code %d", wantedCode, gotCode)
	}

	entriesRouter := &Router{}
	entriesRouter.Handle("/users/[username]/random", func(params map[string]string) http.Handler {
		username := params["username"]
		handler := func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(fmt.Sprintf("random for username: %s", username)))
		}
		return http.HandlerFunc(handler)
	})
	mux.Handle("/api/wordgame/entries/", http.StripPrefix("/api/wordgame/entries", entriesRouter))

	requestToRandom := httptest.NewRequest("GET", "/api/wordgame/entries/users/asmir/random", nil)
	responseFromRandom := httptest.NewRecorder()
	mux.ServeHTTP(responseFromRandom, requestToRandom)
}
