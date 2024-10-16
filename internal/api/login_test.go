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
	"github.com/xavesen/search-api/internal/config"
	"github.com/xavesen/search-api/internal/models"
	"github.com/xavesen/search-api/internal/storage"
	"github.com/xavesen/search-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var loginTests = []struct {
	testName 			string
	payload				models.LoginRequest
	userStorage 		*storage.UserStorageMock
	tokenOp 			*utils.TokenOperatorMock
	expectedCode		int
	expectedResponse 	utils.Response
}{
	{
		testName: "Return 200 and tokens",
		payload: models.LoginRequest{Login: "login", Password: "password"},
		userStorage: &storage.UserStorageMock{
			ExpectedToken: "c0ae478432dcbda4f1b729235ffc4ffe2b211262a62f2345fe2c1378143945b1de77bc53b4a90f8f9832b243956b5ce9a9937bf2e80171a002faeb6014f8abac",
			User: &models.User{Id: "123", Password: "password"},
		},
		tokenOp: &utils.TokenOperatorMock{
			Token: "token",
		},
		expectedCode: 200,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.TokenResponse{
				AccessToken: "token1",
				RefreshToken: "token2",
			},
		},
	},
	{
		testName: "Return 401 when no user with such id",
		payload: models.LoginRequest{Login: "login", Password: "password"},
		userStorage: &storage.UserStorageMock{
			GetUserErr: mongo.ErrNoDocuments,
		},
		tokenOp: &utils.TokenOperatorMock{},
		expectedCode: 401,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Unauthorized",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on db error while getting user info",
		payload: models.LoginRequest{Login: "login", Password: "password"},
		userStorage: &storage.UserStorageMock{
			GetUserErr: errors.New("random error"),
		},
		tokenOp: &utils.TokenOperatorMock{},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 401 with wrong password",
		payload: models.LoginRequest{Login: "login", Password: "password"},
		userStorage: &storage.UserStorageMock{
			User: &models.User{Id: "123", Password: "not_password"},
		},
		tokenOp: &utils.TokenOperatorMock{},
		expectedCode: 401,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Unauthorized",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on error generating token",
		payload: models.LoginRequest{Login: "login", Password: "password"},
		userStorage: &storage.UserStorageMock{
			User: &models.User{Id: "123", Password: "password"},
		},
		tokenOp: &utils.TokenOperatorMock{
			GenerateErr: errors.New("random error"),
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
	{
		testName: "Return 500 on db error while setting refresh token",
		payload: models.LoginRequest{Login: "login", Password: "password"},
		userStorage: &storage.UserStorageMock{
			ExpectedToken: "c0ae478432dcbda4f1b729235ffc4ffe2b211262a62f2345fe2c1378143945b1de77bc53b4a90f8f9832b243956b5ce9a9937bf2e80171a002faeb6014f8abac",
			User: &models.User{Id: "123", Password: "password"},
			GetUserErr: errors.New("random error"),
		},
		tokenOp: &utils.TokenOperatorMock{
			Token: "token",
		},
		expectedCode: 500,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
}

func TestLoginHandler(t *testing.T) {
	config := &config.Config{
		JwtKey: []byte("aaa"),
		TokenHeaderName: "aaa",
		JwtSalt: "aaa",
		JwtAccessTTL: 1,
		JwtRefreshTTL: 2,
	}
	for i, test := range loginTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", nil, nil, test.userStorage, config, test.tokenOp)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaledPayload))
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