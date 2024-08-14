package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaymentReturn_APPROVED(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new server instance
	server := &Server{}

	// Create a test HTTP request
	req, err := http.NewRequest(http.MethodGet, "/payment-return?status=APPROVED&orderID=someOID&merchantOrderID=someMOID", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	resp := httptest.NewRecorder()

	// Setup router and route
	r := gin.Default()
	r.GET("/payment-return", server.paymentReturn)

	// Serve the request
	r.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal the response body into a PaymentReturn struct
	var result PaymentReturn
	err = json.Unmarshal([]byte(resp.Body.String()), &result)
	assert.NoError(t, err)

	// Validate the returned JSON
	assert.Equal(t, "APPROVED", result.Status)
	assert.Equal(t, "someOID", result.OrderID)
	assert.Equal(t, "someMOID", result.MerchantOrderID)
	assert.Equal(t, UserMessage, result.UserMessage)
	assert.Equal(t, "", result.ErrorMessage)
	assert.Equal(t, "", result.PossibleStatuses)
}

func TestPaymentReturn_ErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new server instance
	server := &Server{}

	// Create a test HTTP request with an error message
	req, err := http.NewRequest(http.MethodGet, "/payment-return?status=ERROR&orderID=someOID&merchantOrderID=someMOID&errorMessage=SomeError", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	resp := httptest.NewRecorder()

	// Setup router and route
	r := gin.Default()
	r.GET("/payment-return", server.paymentReturn)

	// Serve the request
	r.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal the response body into a PaymentReturn struct
	var result PaymentReturn
	err = json.Unmarshal([]byte(resp.Body.String()), &result)
	assert.NoError(t, err)

	// Validate the returned JSON
	assert.Equal(t, "ERROR", result.Status)
	assert.Equal(t, "someOID", result.OrderID)
	assert.Equal(t, "someMOID", result.MerchantOrderID)
	assert.Equal(t, UserMessage, result.UserMessage)
	assert.Equal(t, "SomeError", result.ErrorMessage)
	assert.Equal(t, PossibleStatusesLink, result.PossibleStatuses)
}

func TestPaymentReturn_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new server instance
	server := &Server{}

	// Create a test HTTP request missing required fields
	req, err := http.NewRequest(http.MethodGet, "/payment-return?status=ERROR", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	resp := httptest.NewRecorder()

	// Setup router and route
	r := gin.Default()
	r.GET("/payment-return", server.paymentReturn)

	// Serve the request
	r.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
