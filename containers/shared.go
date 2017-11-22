package containers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var serviceRegistryURL string = "http://[::1]:4646"

type Message struct {
	Number uint64 `json:"number"`
	Text   string `json:"text"`
}

type Service struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

func RegisterService(service *Service) error {
	serviceJSON, err := json.Marshal(service)
	if err != nil {
		return err
	}

	resp, err := http.Post(serviceRegistryURL+"/register", "application/json", bytes.NewBuffer(serviceJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error registering service %s: status code %s", service, resp.StatusCode)
	}
	return nil
}

func RetrieveService(serviceName string) ([]Service, error) {
	resp, err := http.Post(serviceRegistryURL+"/retrieve", "application/json", bytes.NewBufferString(serviceName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error retrieving service %s: status code %s", serviceName, resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var services []Service
	err = json.Unmarshal(respBody, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}
