package client

import (
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

type Client struct {
    client *http.Client
}

func (httpClient *Client) Get(url string) (string, error) {

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

func NewClient() *Client {
    httpClient := new(Client)
    httpClient.client = &http.Client{
        Timeout: time.Second * 5,
    }

    return httpClient
}

type HttpClientBuilder func() HttpClient

type HttpClient interface {
    Get(address string) (string, error)
}

