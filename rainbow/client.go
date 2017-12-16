package rainbow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	serverAddr string
}

func NewClient(serverAddr string) *Client {
	return &Client{serverAddr: serverAddr}
}

func (c *Client) List() ([]Deployment, error) {
	resp, err := http.Get(c.serverAddr + "/deployments")
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var deployments []Deployment
	err = json.Unmarshal(responseBody, &deployments)
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (c *Client) Create(deployment Deployment) (*Deployment, error) {
	body, err := json.Marshal(deployment)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(c.serverAddr+"/deployments", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var createdDeployment Deployment
	err = json.Unmarshal(responseBody, &createdDeployment)
	if err != nil {
		return nil, err
	}
	return &createdDeployment, nil
}

func (c *Client) Get(deploymentName string) (*Deployment, error) {
	resp, err := http.Get(c.serverAddr + "/deployments/" + deploymentName)
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var deployment Deployment
	err = json.Unmarshal(responseBody, &deployment)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (c *Client) Delete(deploymentName string) error {
	return fmt.Errorf("DELETE NOT YET IMPLEMENTED")
}
