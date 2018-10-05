package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	client  *http.Client
	limiter <-chan time.Time
}

func (httpClient *Client) Get(url string) (string, error) {
	<-httpClient.limiter
	resp, err := httpClient.client.Get(url)
	if err != nil {
		log.Printf("Unable to Get from %s: %s\n", url, err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read body from %s: %s\n", url, err)
		return "", err
	}
	return string(body), nil
}

func NewClientBuilder() func(<-chan time.Time) HttpClient {
	return func(limiter <-chan time.Time) HttpClient {
		httpClient := new(Client)
		httpClient.client = &http.Client{
			Timeout: time.Second * 5,
		}
		httpClient.limiter = limiter

		return httpClient
	}
}

type HttpClientBuilder func(<-chan time.Time) HttpClient

type HttpClient interface {
	Get(address string) (string, error)
}
