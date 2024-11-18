package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SearchResult struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type Fields struct {
	Summary string `json:"summary"`
	Status  struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"status"`
	Description string `json:"description"`
}

type searchPayload struct {
	Jql    string `json:"jql"`
	Fields string `json:"fields"`
}

type SearchOption func(p *searchPayload)

func WithSearchFields(fields []string) SearchOption {
	return func(p *searchPayload) {
		p.Fields = strings.Join(fields, ",")
	}
}

func (c *Client) Search(jql string, opts ...SearchOption) (*SearchResult, error) {

	payload := searchPayload{
		Jql:    jql,
		Fields: "*navigable",
	}

	for _, opt := range opts {
		opt(&payload)
	}

	req, err := http.NewRequest("GET", "/rest/api/2/search", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("jql", jql)
	q.Add("fields", payload.Fields)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()
	result := &SearchResult{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
