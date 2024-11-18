package jira

import (
	"net/http"
	"net/url"
)

type Client struct {
	h               *http.Client
	username, token string
	baseUrl         *url.URL
}

type NewClientOptions func(*Client)

func WithApiToken(token string) NewClientOptions {
	return func(c *Client) {
		c.token = token
	}
}

func WithUsername(username string) NewClientOptions {
	return func(c *Client) {
		c.username = username
	}
}

func WithBaseUrl(baseUrl string) NewClientOptions {
	return func(c *Client) {
		u, err := url.Parse(baseUrl)
		if err != nil {
			panic(err)
		}
		c.baseUrl = u
	}
}

func NewClient(h *http.Client, opts ...NewClientOptions) *Client {
	if h == nil {
		h = http.DefaultClient
	}
	c := &Client{h: h}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.URL = c.baseUrl.ResolveReference(req.URL)
	req.Header.Set("Authorization", "Basic "+c.token)
	req.SetBasicAuth(c.username, c.token)
	return c.h.Do(req)
}
