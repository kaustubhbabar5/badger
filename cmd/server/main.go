package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaustubhbabar5/badger/badger"
)

func main() {
	badge := badger.NewBadge(100, time.Second*10)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		badge.Push(string(body))
	}).Methods("POST")

	http.ListenAndServe("127.0.0.1:8000", r)
}
