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

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/magiconair/properties/assert"
	"github.com/xavesen/search-api/internal/config"
	"github.com/xavesen/search-api/internal/models"
	"github.com/xavesen/search-api/internal/queue"
	"github.com/xavesen/search-api/internal/storage"
	"github.com/xavesen/search-api/internal/utils"
)

var indexDocumentsTests = []struct {
	testName 			string
	docStorage 			*storage.DocStorageMock
	queue				*queue.QueueMock
	userStorage 		*storage.UserStorageMock
	payload				*models.DocumentsForIndexing
	tokenOp 			*utils.TokenOperatorMock
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 403 when user doesn't have access to index",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: true,
		},
		queue: &queue.QueueMock{
			Error: nil,
		},
		userStorage: &storage.UserStorageMock{
			IndexAccess: false,
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 403,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Index doesn't exist or you don't have access to it",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on error while checking if user has access to index",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: true,
		},
		queue: &queue.QueueMock{
			Error: nil,
		},
		userStorage: &storage.UserStorageMock{
			IndexAccess: false,
			IndexRightsError: errors.New("random error"),
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 401 with invalid token",
		docStorage: &storage.DocStorageMock{
			EsIndexExists: true,
		},
		queue: &queue.QueueMock{
			Error: nil,
		},
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
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
		tokenOp: &utils.TokenOperatorMock{TokenValid: false},
		expectedCode: 401,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Unauthorized",
			Data: nil,
		},
	},
}

func TestIndexDocumentsHandler(t *testing.T) {
	config := &config.Config{
		JwtKey: []byte("aaa"),
		TokenHeaderName: "aaa",
		JwtSalt: "aaa",
	}
	for i, test := range indexDocumentsTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.queue, test.docStorage, test.userStorage, config, test.tokenOp)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/indexDocuments", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}
		req.Header.Add(config.TokenHeaderName, "aaa")

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
	userStorage 		*storage.UserStorageMock
	payload				*models.DocumentSearchRequest
	tokenOp 			*utils.TokenOperatorMock
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 403 if user doesn't have access to index",
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: false,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 403,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Index doesn't exist or you don't have access to it",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on error while checking if user has access to index",
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: false,
			IndexRightsError: errors.New("random error"),
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 401 with invalid token",
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
		userStorage: &storage.UserStorageMock{
			IndexAccess: true,
		},
		payload: &models.DocumentSearchRequest{
			Index: "test",
			Query: "search",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: false},
		expectedCode: 401,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Unauthorized",
			Data: nil,
		},
	},
}

func TestSearchDocumentsHandler(t *testing.T) {
	config := &config.Config{
		JwtKey: []byte("aaa"),
		TokenHeaderName: "aaa",
		JwtSalt: "aaa",
	}
	for i, test := range searchDocumentsTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", nil, test.docStorage, test.userStorage, config, test.tokenOp)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/searchDocuments", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}

		req.Header.Add(config.TokenHeaderName, "aaa")

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
	userStorage 		*storage.UserStorageMock
	payload				*models.CreateIndexRequest
	tokenOp 			*utils.TokenOperatorMock
	expectedCode		int
	expectedResponse 	utils.Response
}{
	{
		testName: "Return 200",
		docStorage: &storage.DocStorageMock{},
		userStorage: &storage.UserStorageMock{},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
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
		userStorage: &storage.UserStorageMock{},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 409,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Index with such name already exists",
			Data: nil,
		},
	},
	{
		testName: "Return 500 when user storage returns an error",
		docStorage: &storage.DocStorageMock{},
		userStorage: &storage.UserStorageMock{
			AddIndexError: errors.New("random error"),
		},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: true},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 401 with invalid token",
		docStorage: &storage.DocStorageMock{},
		userStorage: &storage.UserStorageMock{},
		payload: &models.CreateIndexRequest{
			Index: "test",
		},
		tokenOp: &utils.TokenOperatorMock{TokenValid: false},
		expectedCode: 401,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Unauthorized",
			Data: nil,
		},
	},
}

func TestCreateIndexHandler(t *testing.T) {
	config := &config.Config{
		JwtKey: []byte("aaa"),
		TokenHeaderName: "aaa",
		JwtSalt: "aaa",
	}
	for i, test := range createIndexHandlerTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", nil, test.docStorage, test.userStorage, config, test.tokenOp)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/createIndex", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}

		req.Header.Add(config.TokenHeaderName, "aaa")

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