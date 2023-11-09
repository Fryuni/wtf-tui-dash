package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/wtfutil/wtf/logger"
	"io"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	client     *retryablehttp.Client
	userId     string
	apiToken   string
	clientName string
}

type apiResponse[T any] struct {
	Success    bool   `json:"success"`
	Data       T      `json:"data"`
	AppVersion string `json:"appVersion"`
}

func NewClient(userId, apiToken string) *Client {
	c := &Client{
		client:     retryablehttp.NewClient(),
		userId:     userId,
		apiToken:   apiToken,
		clientName: fmt.Sprintf("%s-WTFDash", userId),
	}

	c.client.Logger = log.Default()

	return c
}

const (
	apiUserHeader  = "X-Api-User"
	apiTokenHeader = "X-Api-Key"
	clientHeader   = "X-Client"

	baseApiURL = "https://habitica.com/api/v3/"
)

func (c *Client) doRequest(method, path string, body, response any) error {
	if body != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return fmt.Errorf("json encoding body: %w", err)
		}
		body = buf
	}

	req, err := retryablehttp.NewRequest(method, baseApiURL+strings.TrimPrefix(path, "/"), body)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	req.Header.Set(apiUserHeader, c.userId)
	req.Header.Set(apiTokenHeader, c.apiToken)
	req.Header.Set(clientHeader, c.clientName)

	logger.LogJson("Request headers", req.Header)

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if res.StatusCode >= http.StatusBadRequest {
		response, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		return fmt.Errorf("API error: %s", response)
	}

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}

	return nil
}
