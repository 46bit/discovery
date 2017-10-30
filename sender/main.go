package main

import (
	"bytes"
	"encoding/json"
	"github.com/46bit/discovery"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.Println("Started...")

	url := "http://localhost:4700/"
	log.Printf("Sending to %s\n", url)

	var i uint64
	for i = 0; true; i++ {
		message := discovery.Message{
			Number: i,
			Text:   "Sender sending.",
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

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Response: %s\n", string(body))
	}
}
