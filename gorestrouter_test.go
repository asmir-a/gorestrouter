package gorestrouter

import (
	"encoding/json"
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
