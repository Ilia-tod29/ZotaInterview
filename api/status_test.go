package api

import (
	"ZotaInterview/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mocking the server's configuration
var mockConfig = util.Config{
	HTTPServerAddress: "0.0.0.0:8080",
	SecretKey:         "someSecretkey",
	EndpointId:        "someEndpoint",
	MerchantId:        "someMerchant",
}

// TestCheckDepositStatus_APPROVED tests the successful handling of the checkDepositStatus endpoint.
func TestCheckDepositStatus_APPROVED(t *testing.T) {
	// Mock ZotaClient
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	server := &Server{
		router:     router,
		zotaClient: &MockZotaClient{},
	}
	router.GET("/status", server.checkDepositStatus)

	successfulRespBody, _ := json.Marshal(map[string]interface{}{
		"Status": "APPROVED",
	})

	tests := []struct {
		name           string
		qparams        string
		mockResponse   *http.Response
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:    "Successful request (APPROVED)",
			qparams: "orderID=someOID&merchantOrderID=someMOID",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(successfulRespBody)),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"Status": "APPROVED",
			},
		},
		{
			name:           "Not all request params provided",
			qparams:        "orderID=someOID",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Key: 'StatusRequest.MerchantOrderID' Error:Field validation for 'MerchantOrderID' failed on the 'required' tag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockClient := server.zotaClient.(*MockZotaClient)
			mockClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			// Create a request and response recorder
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/status?%s", tt.qparams), nil)
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check response
			require.Equal(t, tt.expectedStatus, resp.Code)
			var respBody map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &respBody)
			require.Equal(t, tt.expectedBody, respBody)
		})
	}
}

// TestGenerateStatusFields tests the generation of signature and timestamp in StatusRequest.
func TestGenerateStatusFields(t *testing.T) {
	server := &Server{config: mockConfig}
	req := &StatusRequest{
		OrderID:         "someOID",
		MerchantOrderID: "someMOID",
	}

	server.generateStatusFields(req)

	require.Equal(t, "someMerchant", req.MerchantID)
	require.NotEmpty(t, req.Timestamp)
	require.NotEmpty(t, req.Signature)

	expectedSignature := util.GenerateSignature("someMerchant" + "someMOID" + "someOID" + req.Timestamp + "someSecretkey")
	require.Equal(t, expectedSignature, req.Signature)
}
