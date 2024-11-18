package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetBoardById(boardId string) (*Board, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/rest/agile/1.0/board/%s", boardId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	b := &Board{}
	err = decoder.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return b, nil
}

func (c *Client) GetBoards() (*BoardList, error) {
	req, err := http.NewRequest("GET", "/rest/agile/1.0/board", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	b := &BoardList{}
	err = decoder.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return b, nil
}

type Board struct {
	Id       int `json:"id"`
	Location struct {
		DisplayName string `json:"displayName"`
		ProjectName string `json:"projectName"`
		ProjectKey  string `json:"projectKey"`
	} `json:"location"`
	Name string `json:"name"`
}

type BoardList struct {
	Values     []Board `json:"values"`
	IsLast     bool    `json:"isLast"`
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
}

type BoardConfig struct {
	ColumnConfig struct {
		Columns []BoardColumn `json:"columns"`
	} `json:"columnConfig"`
}

type BoardColumn struct {
	Name     string        `json:"name"`
	Statuses []BoardStatus `json:"statuses"`
}

type BoardStatus struct {
	Id string `json:"id"`
}

func (c *Client) GetBoardConfig(boardId string) (*BoardConfig, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/rest/agile/1.0/board/%s/configuration", boardId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	b := &BoardConfig{}
	err = decoder.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return b, nil
}

type Sprint struct {
	Id    int    `json:"id"`
	State string `json:"state"`
	Name  string `json:"name"`
	Goal  string `json:"goal"`
}

type SprintList struct {
	Values     []Sprint `json:"values"`
	IsLast     bool     `json:"isLast"`
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
}

func (c *Client) GetBoardActiveSprints(boardId string) (*SprintList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/rest/agile/1.0/board/%s/sprint?state=active", boardId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	s := &SprintList{}
	err = decoder.Decode(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s, nil
}
