package api

import (
	"net/http"
	"time"
)

type (
	Task struct {
		Id          string              `json:"id"`
		Up          bool                `json:"up,omitempty"`
		Down        bool                `json:"down,omitempty"`
		Type        string              `json:"type"`
		Notes       string              `json:"notes"`
		Tags        []string            `json:"tags"`
		Value       float64             `json:"value"`
		Priority    float64             `json:"priority"`
		Attribute   string              `json:"attribute"`
		ByHabitica  bool                `json:"byHabitica"`
		CreatedAt   time.Time           `json:"createdAt"`
		UpdatedAt   time.Time           `json:"updatedAt"`
		Text        string              `json:"text"`
		YesterDaily bool                `json:"yesterDaily,omitempty"`
		Completed   bool                `json:"completed,omitempty"`
		Checklist   []TaskChecklistItem `json:"checklist,omitempty"`
		IsDue       bool                `json:"isDue,omitempty"`
	}
	TaskChecklistItem struct {
		Completed bool   `json:"completed"`
		Text      string `json:"text"`
		Id        string `json:"id"`
	}
)

func (c *Client) ListTodos() ([]*Task, error) {
	var response apiResponse[[]*Task]

	err := c.doRequest(http.MethodGet, "/tasks/user?type=todos", nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}
