package main

import (
	"encoding/json"
	"github.com/46bit/discovery/containers"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("Started...")

	counts := map[string]uint64{}

	go func() {
		sum := uint64(0)
		for _, n := range counts {
			sum += n
		}
		log.Printf("aggregator count = %d", sum)
		time.Sleep(2 * time.Second)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, `{"error": "Request body was missing."}`, 400)
			log.Printf("Request body was missing.")
			return
		}

		var m containers.Message
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, `{"error": "Request body was not valid JSON."}`, 400)
			log.Printf("Request body was not valid JSON: %s.\n", err)
			return
		}

		http.Error(w, `{"success": "true"}`, 200)

		counts[m.Text] = m.Number
	})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Error listening: %s\n", err)
	}
	listenPort := listener.Addr().(*net.TCPAddr).Port

	err = containers.RegisterService(&containers.Service{
		Name: "aggregator",
		Host: fmt.Sprintf("[::1]:%d", listenPort),
	})
	if err != nil {
		log.Fatalln(err)
	}

	go func(listenPort int) {
		for {
			err = containers.RegisterService(&containers.Service{
				Name: "aggregator",
				Host: fmt.Sprintf("[::1]:%d", listenPort),
			})
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(1 * time.Second)
		}
	}(listenPort)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatalf("Error serving on port %d: %s\n", listenPort, err)
	}
}
