package typesense

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	apiKey = APIKey{
		Value: "apikey",
		Actions: []APIAction{
			DocAllAPIAction,
			CollectionCreateAPIAction,
		},
		ID:          1,
		Collections: []string{"companies"},
		Description: "API key",
	}
)

func TestCreateAPIKey(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(apiKey)
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	resAPIKey, err := client.CreateAPIKey(apiKey.Description, apiKey.Actions, apiKey.Collections)
	if err != nil {
		t.Errorf("Expected to receive nil error, received: %v", err)
	}
	if !reflect.DeepEqual(apiKey, *resAPIKey) {
		t.Errorf("Expected api key %+v, received %+v", apiKey, *resAPIKey)
	}
}

func TestCreateAPIKey_badRequest(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(apiKey)
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	resAPIKey, err := client.CreateAPIKey(apiKey.Description, apiKey.Actions, apiKey.Collections)
	if err != nil {
		t.Errorf("Expected to receive nil error, received: %v", err)
	}
	if !reflect.DeepEqual(apiKey, *resAPIKey) {
		t.Errorf("Expected api key %+v, received %+v", apiKey, *resAPIKey)
	}
}

func TestCreateAPIKey_unauthorized(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateAPIKey(apiKey.Description, apiKey.Actions, apiKey.Collections)
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}

func TestAPIKey(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(apiKey)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.APIKey(apiKey.ID)
	if err != nil {
		t.Errorf("Expected to receive error nil, received %v", err)
	}
}

func TestAPIKey_unauthorized(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.APIKey(apiKey.ID)
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}

func TestAPIKeys(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal([]APIKey{apiKey})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.APIKeys()
	if err != nil {
		t.Errorf("Expected to receive error nil, received %v", err)
	}
}

func TestAPIKeys_unauthorized(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal([]APIKey{apiKey})
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.APIKeys()
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}

func TestDeleteAPIKey(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.DeleteAPIKey(apiKey.ID)
	if err != nil {
		t.Errorf("Expected to receive error nil, received %v", err)
	}
}

func TestDeleteAPIKey_unauthorized(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.DeleteAPIKey(apiKey.ID)
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}
