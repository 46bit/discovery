package main

import (
	"encoding/json"
	"github.com/46bit/discovery/containers"
	"net/http"
	"time"
)

type TimestampedService struct {
	Service   containers.Service
	Timestamp time.Time
}

func main() {
	serviceRegistry := map[string]map[string]TimestampedService{}

	http.HandleFunc("/services/register", func(w http.ResponseWriter, r *http.Request) {
		var service containers.Service
		err := json.NewDecoder(r.Body).Decode(&service)
		if err != nil {
			http.Error(w, `{"error": "Request body was not valid JSON."}`, 400)
			return
		}
		timestampedService := TimestampedService{
			Service:   service,
			Timestamp: time.Now(),
		}

		if _, ok := serviceRegistry[service.Name]; !ok {
			serviceRegistry[service.Name] = map[string]TimestampedService{}
		}
		serviceRegistry[service.Name][service.Host] = timestampedService
	})

	http.HandleFunc("/services/retrieve", func(w http.ResponseWriter, r *http.Request) {
		var serviceName string
		err := json.NewDecoder(r.Body).Decode(&serviceName)
		if err != nil {
			http.Error(w, `{"error": "Request body was not valid JSON."}`, 400)
			return
		}

		if _, ok := serviceRegistry[serviceName]; !ok {
			http.Error(w, `{"error": "Not found."}`, 404)
			return
		}

		services := []containers.Service{}
		for _, timestampedService := range serviceRegistry[serviceName] {
			age := time.Since(timestampedService.Timestamp)
			if age.Seconds() < 5 {
				services = append(services, timestampedService.Service)
			} else {
				delete(serviceRegistry[serviceName], timestampedService.Service.Name)
			}
		}
		if len(serviceRegistry[serviceName]) == 0 {
			delete(serviceRegistry, serviceName)
		}

		servicesJSON, err := json.Marshal(&services)
		if err != nil {
			http.Error(w, `{"error": "Could not marshal services."}`, 500)
			return
		}
		http.Error(w, string(servicesJSON), 200)
	})

	http.ListenAndServe(":4646", nil)
}
