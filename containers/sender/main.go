package main

import (
	"bytes"
	"encoding/json"
	"github.com/46bit/discovery/containers"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Started...")

	var j uint64
	var url string
	for j = 0; j < 100; j++ {
		services, err := containers.RetrieveService("receiver")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("found services %s\n", services)
		if len(services) > 0 {
			url = "http://" + services[0].Host + "/"
			break
		}
		time.Sleep(time.Second)
	}
	if url == "" {
		log.Fatalln("No receiver service found in time.")
	}

	log.Printf("Sending to %s\n", url)

	var i uint64
	for i = 0; true; i++ {
		func(i uint64) {
			message := containers.Message{
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

			//body, err := ioutil.ReadAll(resp.Body)
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//log.Printf("Response: %s\n", string(body))
		}(i)
	}
}
