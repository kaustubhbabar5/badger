package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaustubhbabar5/badger/badger"
)

func main() {
	push := badger.Initiate(100, 10*time.Second)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		push <- string(body)

	}).Methods("POST")

	http.ListenAndServe("127.0.0.1:8000", r)
}
