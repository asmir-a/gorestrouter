package main

import (
	"net/http"
)

func usernameProvider(username string) http.Handler {
	newFunc := func(w http.ResponseWriter, req *http.Request) {
		//uses username
	}
	return http.HandlerFunc(newFunc)
}

func wordsHandler(w http.ResponseWriter, req *http.Request) {
}

func getUsernameHandler(params map[string]string) http.Handler {
	username := params["username"] //we can use something similar to enums as well
	newFunc := func(w http.ResponseWriter, req *http.Request) {
		//use username
		w.Write([]byte(username))
	}
	return http.HandlerFunc(newFunc)
}
