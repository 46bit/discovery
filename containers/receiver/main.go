package main

import (
	"bytes"
	"encoding/json"
	"github.com/46bit/discovery/containers"
	//"github.com/davecgh/go-spew/spew"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {
	log.Println("Started...")

	var j uint64
	var url string
	for j = 0; j < 100; j++ {
		services, err := containers.RetrieveService("aggregator")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("found services %s\n", services)
		if len(services) > 0 {
			serviceBigIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(services))))
			if err != nil {
				log.Fatalln(err)
			}
			serviceIndex := int(serviceBigIndex.Int64())
			url = "http://" + services[serviceIndex].Host + "/"
			break
		}
		time.Sleep(time.Second)
	}
	if url == "" {
		log.Fatalln("No aggregator service found in time.")
	}

	var requests *uint64 = new(uint64)
	atomic.StoreUint64(requests, 0)

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
		//log.Println("Request was valid: %s", spew.Sdump(m))
		requestsValue := atomic.AddUint64(requests, 1)
		if requestsValue%1000 == 0 {
			fmt.Println(requestsValue)
		}
	})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Error listening: %s\n", err)
	}
	listenPort := listener.Addr().(*net.TCPAddr).Port

	err = containers.RegisterService(&containers.Service{
		Name: "receiver",
		Host: fmt.Sprintf("[::1]:%d", listenPort),
	})
	if err != nil {
		log.Fatalln(err)
	}

	go func(listenPort int) {
		for {
			err = containers.RegisterService(&containers.Service{
				Name: "receiver",
				Host: fmt.Sprintf("[::1]:%d", listenPort),
			})
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(1 * time.Second)
		}
	}(listenPort)

	go func(url string) {
		for {
			message := containers.Message{
				Number: *requests,
				Text:   fmt.Sprintf("[::1]:%d", listenPort),
			}

			messageJSON, err := json.Marshal(&message)
			if err != nil {
				log.Fatalln(err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
			if err != nil {
				log.Fatalln(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()

			time.Sleep(time.Second)
		}
	}(url)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatalf("Error serving on port %d: %s\n", listenPort, err)
	}
}
