package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/46bit/discovery/rainbow"
	"io/ioutil"
	"net/http"
)

type Client struct {
	serverAddr string
}

func NewClient(serverAddr string) *Client {
	return &Client{serverAddr: serverAddr}
}

func (c *Client) List() ([]rainbow.Deployment, error) {
	resp, err := http.Get(c.serverAddr + "/deployments")
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var deployments []rainbow.Deployment
	err = json.Unmarshal(responseBody, &deployments)
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (c *Client) Create(deployment rainbow.Deployment) (*rainbow.Deployment, error) {
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
	var createdDeployment rainbow.Deployment
	err = json.Unmarshal(responseBody, &createdDeployment)
	if err != nil {
		return nil, err
	}
	return &createdDeployment, nil
}

func (c *Client) Get(deploymentName string) (*rainbow.Deployment, error) {
	resp, err := http.Get(c.serverAddr + "/deployments/" + deploymentName)
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var deployment rainbow.Deployment
	err = json.Unmarshal(responseBody, &deployment)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (c *Client) Delete(deploymentName string) error {
	return fmt.Errorf("DELETE NOT YET IMPLEMENTED")
}
