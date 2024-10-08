package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/xavesen/search-api/internal/models"
	"github.com/xavesen/search-api/internal/queue"
	"github.com/xavesen/search-api/internal/storage"
	"github.com/xavesen/search-api/internal/utils"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

var indexDocumentsTests = []struct {
	testName 			string
	docStorage 			*storage.DocStorageMock
	queue				*queue.QueueMock
	payload				*models.DocumentsForIndexing
	expectedCode		int
	expectedResponse 	utils.Response
}{
	{
		testName: "Return 200",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: true,
		},
		queue: &queue.QueueMock{
			Error: nil,
		},
		payload: &models.DocumentsForIndexing{
			Index: "test",
			UserId: "1",
			Documents: []models.Document{
				{
					Title: "test",
					Text: "test test test",
				},
				{
					Title: "test1",
					Text: "test1 test1 test1",
				},
			},
		},
		expectedCode: 200,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: nil,
		},
	},
	{
		testName: "Return 500 when doc storage returns an error on index check",
		docStorage: &storage.DocStorageMock{
			IndexError: errors.New("random error"),
			EsIndexExists: false,
		},
		queue: &queue.QueueMock{
			Error: nil,
		},
		payload: &models.DocumentsForIndexing{
			Index: "test",
			UserId: "1",
			Documents: []models.Document{
				{
					Title: "test",
					Text: "test test test",
				},
				{
					Title: "test1",
					Text: "test1 test1 test1",
				},
			},
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 403 when index doesn't exist",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: false,
		},
		queue: &queue.QueueMock{
			Error: nil,
		},
		payload: &models.DocumentsForIndexing{
			Index: "test",
			UserId: "1",
			Documents: []models.Document{
				{
					Title: "test",
					Text: "test test test",
				},
				{
					Title: "test1",
					Text: "test1 test1 test1",
				},
			},
		},
		expectedCode: 403,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Index doesn't exist or you don't have access to it",
			Data: nil,
		},
	},
	{
		testName: "Return 500 when writing message to queue returns an error",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: true,
		},
		queue: &queue.QueueMock{
			Error: errors.New("random error"),
		},
		payload: &models.DocumentsForIndexing{
			Index: "test",
			UserId: "1",
			Documents: []models.Document{
				{
					Title: "test",
					Text: "test test test",
				},
				{
					Title: "test1",
					Text: "test1 test1 test1",
				},
			},
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
}

func TestIndexDocumentsHandler(t *testing.T) {
	for i, test := range indexDocumentsTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.queue, test.docStorage)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/indexDocuments", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		expectedResp, err := json.Marshal(test.expectedResponse)
		if err != nil {
			t.Fatalf("Unable to marshal expected response, error: %s\n", err)
		}

		assert.Equal(t, rr.Code, test.expectedCode, "wrong response code")
		assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), string(expectedResp), "wrong body contents")
	}
}

var searchDocumentsTests = []struct {
	testName 			string
	docStorage 			*storage.DocStorageMock
	payload				*models.DocumentSearchRequest
	expectedCode		int
	expectedResponse 	utils.Response
}{
	{
		testName: "Return 200",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: true,
			Documents: []models.Document{
				{
					Title: "test",
					Text: "test test test",
				},
				{
					Title: "test1",
					Text: "test1 test1 test1",
				},
			},
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		expectedCode: 200,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: []models.Document{
				{
					Title: "test",
					Text: "test test test",
				},
				{
					Title: "test1",
					Text: "test1 test1 test1",
				},
			},
		},
	},
	{
		testName: "Return 500 when doc storage returns an error on index check",
		docStorage: &storage.DocStorageMock{
			IndexError: errors.New("random error"),
			EsIndexExists: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 403 when index doesn't exist",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: false,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		expectedCode: 403,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Index doesn't exist or you don't have access to it",
			Data: nil,
		},
	},
	{
		testName: "Return 500 when search request fails",
		docStorage: &storage.DocStorageMock{
			SearchError: errors.New("random error"),
			EsIndexExists: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
}

func TestSearchDocumentsHandler(t *testing.T) {
	for i, test := range searchDocumentsTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", nil, test.docStorage)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/searchDocuments", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		expectedResp, err := json.Marshal(test.expectedResponse)
		if err != nil {
			t.Fatalf("Unable to marshal expected response, error: %s\n", err)
		}

		assert.Equal(t, rr.Code, test.expectedCode, "wrong response code")
		assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), string(expectedResp), "wrong body contents")
	}
}

var createIndexHandlerTests = []struct {
	testName 			string
	docStorage 			*storage.DocStorageMock
	payload				*models.CreateIndexRequest
	expectedCode		int
	expectedResponse 	utils.Response
}{
	{
		testName: "Return 200",
		docStorage: &storage.DocStorageMock{},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		expectedCode: 200,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on random error creating index",
		docStorage: &storage.DocStorageMock{
			CreateError: errors.New("random error"),
		},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on ElasticSearchError (except resource already exists)",
		docStorage: &storage.DocStorageMock{
			CreateError: types.NewElasticsearchError(),
		},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 409 on ElasticSearchError for resource already exists",
		docStorage: &storage.DocStorageMock{
			CreateError: &types.ElasticsearchError{Status: 400, ErrorCause: types.ErrorCause{Type: storage.ErrResourceAlreadyExists}},
		},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		expectedCode: 409,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Index with such name already exists",
			Data: nil,
		},
	},
}

func TestCreateIndexHandler(t *testing.T) {
	for i, test := range createIndexHandlerTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", nil, test.docStorage)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/createIndex", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		expectedResp, err := json.Marshal(test.expectedResponse)
		if err != nil {
			t.Fatalf("Unable to marshal expected response, error: %s\n", err)
		}

		assert.Equal(t, rr.Code, test.expectedCode, "wrong response code")
		assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), string(expectedResp), "wrong body contents")
	}
}