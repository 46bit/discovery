package main

import (
	"encoding/json"
	"github.com/46bit/discovery/containers"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	log.Println("Started...")

	var requests *uint64 = new(uint64)
	atomic.StoreUint64(requests, 0)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, `{"error": "Request body was missing."}`, 400)
			log.Println("Request body was missing.")
			return
		}

		var m containers.Message
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, `{"error": "Request body was not valid JSON."}`, 400)
			log.Println("Request body was not valid JSON: %s.\n", err)
			return
		}

		http.Error(w, `{"success": "true"}`, 200)
		//log.Println("Request was valid: %s", spew.Sdump(m))
		requestsValue := atomic.AddUint64(requests, 1)
		if requestsValue%1000 == 0 {
			fmt.Println(requestsValue)
		}
	})

	http.ListenAndServe(":4700", nil)
}
