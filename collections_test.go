package typesense

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var (
	testCollectionSchema = CollectionSchema{
		Name: "companies",
		Fields: []CollectionField{
			{
				Name: "name",
				Type: "string",
			},
		},
	}
	testCollection = Collection{
		testCollectionSchema,
		0,
		0,
	}
	testOverride = Override{
		ID: "customize-apple",
		Rule: OverrideRule{
			Match: "exact",
			Query: "apple",
		},
	}
)

func TestCollectionField(t *testing.T) {
	_ = CollectionField{
		Name:  "test",
		Type:  "string",
		Facet: true,
	}
}

func TestCreateCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionData, _ := json.Marshal(testCollection)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(collectionData)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collectionResp, err := client.CreateCollection(testCollectionSchema)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if collectionResp == nil {
		t.Errorf("Expected to receive a collection as the response, received %v", collectionResp)
	}
}

func TestCreateCollection_nameRequired(t *testing.T) {
	testData := CollectionSchema{Fields: []CollectionField{{Name: "field", Type: "string"}}}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionData, _ := json.Marshal(&testData)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(collectionData)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateCollection(testData)
	if err != ErrCollectionNameRequired {
		t.Errorf("Expected to receive error %v, received error %v", ErrCollectionNameRequired, err)
	}
}

func TestCreateCollection_fieldsRequired(t *testing.T) {
	testData := CollectionSchema{Name: "new-collection"}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionData, _ := json.Marshal(&testData)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(collectionData)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateCollection(testData)
	if err != ErrCollectionFieldsRequired {
		t.Errorf("Expected to receive error %v, received error %v", ErrCollectionFieldsRequired, err)
	}
}

func TestCreateCollection_conflict(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusConflict,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "collection already exists"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateCollection(testCollectionSchema)
	if err != ErrCollectionDuplicate {
		t.Errorf("Expected to receive error message %v, received error %v", ErrCollectionDuplicate, err)
	}
}

func TestRetrieveCollections(t *testing.T) {
	jsonBody := `[{"name": "companies", "num_documents": 0, "fields": [{"name": "name", "type": "string", "facet": false}]}]`
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(jsonBody)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collections, err := client.RetrieveCollections()
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if collections == nil {
		t.Errorf("Expected to receive collections, received nil")
	}
}

func TestRetrieveCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionJSON, _ := json.Marshal(&testCollection)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(collectionJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collection, err := client.RetrieveCollection(testCollection.Name)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if !reflect.DeepEqual(*collection, testCollection) {
		t.Errorf("Expected to receive %v, received %v", testCollection, *collection)
	}
}

func TestRetrieveCollection_notFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"json": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.RetrieveCollection(testCollection.Name)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestDeleteCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionJSON, _ := json.Marshal(&testCollection)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(collectionJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collection, err := client.DeleteCollection(testCollection.Name)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if !reflect.DeepEqual(*collection, testCollection) {
		t.Errorf("Expected to receive %v, received %v", testCollection, *collection)
	}
}

func TestDeleteCollection_notFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"json": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.DeleteCollection(testCollection.Name)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestOverrideCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		reqBody, _ := ioutil.ReadAll(req.Body)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(reqBody)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.OverrideCollection(testCollection.Name, testOverride)
	if err != nil {
		t.Errorf("Expected to receive nil error, received %v", err)
	}
}

func TestOverrideCollection_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.OverrideCollection(testCollection.Name, testOverride)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestOverrideCollection_unauthorized(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.OverrideCollection(testCollection.Name, testOverride)
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}

func TestRetrieveOverrides(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		bodyData := make(map[string][]*Override)
		bodyData["overrides"] = []*Override{&testOverride}
		body, _ := json.Marshal(bodyData)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	overrides, err := client.RetrieveOverrides(testCollection.Name)
	if err != nil {
		t.Errorf("Expected to receive nil error, received %v", err)
	}
	if len(overrides) == 0 {
		t.Errorf("Expected to receive at least one override, received 0")
	}
}

func TestRetrieveOverrides_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.RetrieveOverrides(testCollection.Name)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestRetrieveOverrides_unauthorized(t *testing.T) {
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
	_, err := client.RetrieveOverrides(testCollection.Name)
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}

func TestDeleteOverride(t *testing.T) {
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
	err := client.DeleteOverride(testCollection.Name, testOverride.ID)
	if err != nil {
		t.Errorf("Expected to receive error nil, received %v", err)
	}
}

func TestDeleteOverride_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.DeleteOverride(testCollection.Name, testOverride.ID)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestDeleteOverride_unauthorized(t *testing.T) {
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
	err := client.DeleteOverride(testCollection.Name, testOverride.ID)
	if err != ErrUnauthorized {
		t.Errorf("Expected to receive error %v, received %v", ErrUnauthorized, err)
	}
}
