package main

import (
	"encoding/json"
	"github.com/46bit/discovery"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
)

func main() {
	log.Println("Started...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, `{"error": "Request body was missing."}`, 400)
			log.Println("Request body was missing.")
			return
		}

		var m discovery.Message
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, `{"error": "Request body was not valid JSON."}`, 400)
			log.Println("Request body was not valid JSON: %s.\n", err)
			return
		}

		http.Error(w, `{"success": "true"}`, 200)
		log.Println("Request was valid: %s", spew.Sdump(m))
	})

	http.ListenAndServe(":4700", nil)
}
