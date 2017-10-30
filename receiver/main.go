package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
)

type Message struct {
	number uint64
	text   string
}

func main() {
	log.Println("Started...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, `{"error": "Request body was missing."}`, 400)
			log.Println("Request body was missing.")
			return
		}

		var m Message
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), 400)
			log.Println("Request body was not valid JSON.")
			return
		}

		log.Println("Request was valid: %s", spew.Sdump(m))
	})

	http.ListenAndServe(":4700", nil)
}
