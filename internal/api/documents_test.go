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
			Error: nil,
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
			Error: errors.New("random error"),
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
			Error: nil,
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
			Error: nil,
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