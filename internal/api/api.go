package api

import (
	"errors"
	"time"
)

type Client struct {
	token    string
	host     string
	endpoint string
}

var (
	ErrorNotFound            error         = errors.New("not found")
	clientTimeout            time.Duration = 60
	rateLimitingNumOfRetries int           = 3
)

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

func (c *Client) GetToken() string {
	return c.token
}

func (c *Client) GetHost() string {
	return c.host
}

func (c *Client) GetEndpoint() string {
	return c.endpoint
}
