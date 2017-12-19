package rainbow

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	serverAddr string
}

func NewClient(serverAddr string) *Client {
	return &Client{serverAddr: serverAddr}
}

func (c *Client) ListDeployments() ([]Deployment, error) {
	resp, err := httpClient().Get(c.serverAddr + "/deployments")
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

func (c *Client) CreateDeployment(deployment Deployment) error {
	body, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	_, err = httpClient().Post(c.serverAddr+"/deployments", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetDeployment(name string) (*Deployment, error) {
	resp, err := httpClient().Get(c.serverAddr + "/deployments/" + url.PathEscape(name))
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

func (c *Client) DeleteDeployment(name string) error {
	req, err := http.NewRequest("DELETE", c.serverAddr+"/deployments/"+url.PathEscape(name), nil)
	if err != nil {
		return err
	}
	_, err = httpClient().Do(req)
	if err != nil {
		return err
	}
	return nil
}

// @TODO: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
func httpClient() *http.Client {
	return &http.Client{}
}
