package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Unit test for depositMoney function
func TestDepositMoney(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	server := &Server{
		router:     router,
		zotaClient: &MockZotaClient{},
	}

	// Define the route
	router.POST("/deposit", server.depositMoney)

	// Define a successful response body
	successResponseBody := map[string]interface{}{
		"code": "200",
		"data": map[string]interface{}{
			"depositUrl":      "someUrl",
			"merchantOrderID": "someMOID",
			"orderID":         "someOID",
		},
	}
	// Convert response body to JSON
	responseBody, _ := json.Marshal(successResponseBody)

	validInputBody := CreateDepositRequest{
		MerchantOrderDesc:   "Test order",
		OrderAmount:         "500.00",
		OrderCurrency:       "USD",
		CustomerEmail:       "someEmail@gmail.com",
		CustomerFirstName:   "Some",
		CustomerLastName:    "Name",
		CustomerAddress:     "Some address",
		CustomerCountryCode: "BG",
		CustomerCity:        "Some City",
		CustomerZipCode:     "1234",
		CustomerPhone:       "+359884324571",
		CustomerIP:          "103.106.8.104",
		RedirectUrl:         "http://localhost:8080/some-return-endpoint/",
		CallbackUrl:         "http://localhost:8080/some-callback-url/",
	}

	// Test cases
	tests := []struct {
		name           string
		inputBody      interface{}
		mockResponse   *http.Response
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Successful request",
			inputBody: validInputBody,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code": "200",
				"data": map[string]interface{}{
					"depositUrl":      "someUrl",
					"merchantOrderID": "someMOID",
					"orderID":         "someOID",
				},
			},
		},
		{
			name:           "Invalid request body",
			inputBody:      "{invalid_json}",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "json: cannot unmarshal string into Go value of type api.CreateDepositRequest"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockClient := server.zotaClient.(*MockZotaClient)
			mockClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			// Make request
			body, _ := json.Marshal(tt.inputBody)
			req, _ := http.NewRequest(http.MethodPost, "/deposit", bytes.NewReader(body))
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			// Check response
			require.Equal(t, tt.expectedStatus, resp.Code)
			var respBody map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &respBody)
			require.Equal(t, tt.expectedBody, respBody)
		})
	}
}

// Unit test for validateDepositRequest function
func TestValidateDepositRequest(t *testing.T) {
	server := &Server{}

	// Test cases
	tests := []struct {
		name          string
		inputRequest  CreateDepositRequest
		expectedError error
	}{
		{
			name: "Valid request",
			inputRequest: CreateDepositRequest{
				OrderAmount:     "100.00",
				CustomerZipCode: "12345",
				// provide all other required fields
			},
			expectedError: nil,
		},
		{
			name: "Invalid order amount",
			inputRequest: CreateDepositRequest{
				OrderAmount:     "invalid_amount",
				CustomerZipCode: "12345",
				// provide all other required fields
			},
			expectedError: errors.New("invalid syntax"),
		},
		{
			name: "Invalid zip code",
			inputRequest: CreateDepositRequest{
				OrderAmount:     "100.00",
				CustomerZipCode: "invalid_zip",
				// provide all other required fields
			},
			expectedError: errors.New("invalid syntax"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.validateDepositRequest(&tt.inputRequest)
			if tt.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError.Error())
			}
		})
	}
}
