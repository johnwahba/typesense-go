package typesense

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// DocSearchAPIAction enable the API key to search for documents.
	DocSearchAPIAction = "documents:search"

	// DocGetAPIAction enable the API key to get documents by id.
	DocGetAPIAction = "documents:get"

	// DocDeleteAPIAction enable the API key to delete documents.
	DocDeleteAPIAction = "documents:delete"

	// DocCreateAPIAction enable the API key to create/index new documents.
	DocCreateAPIAction = "documents:create"

	// DocAllAPIAction enable the API key to perform all actions with documents.
	DocAllAPIAction = "documents:*"

	// CollectionGetAPIAction enable the API key to get collections.
	CollectionGetAPIAction = "collections:get"

	// CollectionDeleteAPIAction enable the API key to delete collections.
	CollectionDeleteAPIAction = "collections:delete"

	// CollectionCreateAPIAction enable the API key to create collections.
	CollectionCreateAPIAction = "collections:create"

	// CollectionAllAPIAction enable the API key to perform all actions with collections.
	CollectionAllAPIAction = "collections:*"

	// AllAPIAction enable the API key to perform all actions with documents and collections.
	AllAPIAction = "*"
)

// APIAction represent an action that the APIKey can perform.
type APIAction string

// APIKey represent a Typesense API key.
type APIKey struct {
	Actions     []APIAction `json:"actions"`
	Collections []string    `json:"collections"`
	Description string      `json:"description"`
	ID          int         `json:"id"`
	Value       string      `json:"value"`
}

// CreateAPIKey creates a new API key using the given actions, and collections access.
func (c *Client) CreateAPIKey(description string, actions []APIAction, collections []string) (*APIKey, error) {
	data := map[string]interface{}{
		"description": description,
		"actions":     actions,
		"collections": collections,
	}
	body, _ := json.Marshal(data)
	url := fmt.Sprintf(
		"%s://%s:%s/keys",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
	)
	res, err := c.apiCall(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusBadRequest {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, APIError{Message: "bad request"}
		}
		return nil, apiError
	} else if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	var apiKey APIKey
	if err := json.NewDecoder(res.Body).Decode(&apiKey); err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// APIKey retrieve an API key by id.
func (c *Client) APIKey(id int) (*APIKey, error) {
	url := fmt.Sprintf(
		"%s://%s:%s/keys/%d",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		id,
	)
	res, err := c.apiCall(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	var apiKey APIKey
	if err := json.NewDecoder(res.Body).Decode(&apiKey); err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// APIKeys retrieve metadata about all API keys.
func (c *Client) APIKeys() ([]*APIKey, error) {
	url := fmt.Sprintf(
		"%s://%s:%s/keys",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
	)
	res, err := c.apiCall(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	var apiKeys []*APIKey
	if err := json.NewDecoder(res.Body).Decode(&apiKeys); err != nil {
		return nil, err
	}
	return apiKeys, nil
}

// DeleteAPIKey delete an API key by id.
func (c *Client) DeleteAPIKey(id int) error {
	url := fmt.Sprintf(
		"%s://%s:%s/keys/%d",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		id,
	)
	res, err := c.apiCall(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	return nil
}

// GenerateScopedSearchKey generate a scoped search key that have
// embedded search parameters in it. The options parameter will
// specify the search parameters.
func (c *Client) GenerateScopedSearchKey(searchKey string, options map[string]string) string {
	h := hmac.New(sha256.New, []byte(searchKey))
	j, _ := json.Marshal(options)
	_, err := h.Write(j)
	if err != nil {
		return ""
	}
	keyPrefix := []byte(searchKey)[:4]
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	rawScopedKey := fmt.Sprintf(digest + string(keyPrefix) + string(j))
	return base64.StdEncoding.EncodeToString([]byte(rawScopedKey))
}
