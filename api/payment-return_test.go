package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPaymentReturn_APPROVED tests the data displayed after a successful deposit
func TestPaymentReturn_APPROVED(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{}

	req, err := http.NewRequest(http.MethodGet, "/payment-return?status=APPROVED&orderID=someOID&merchantOrderID=someMOID", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/payment-return", server.paymentReturn)

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var result PaymentReturn
	err = json.Unmarshal([]byte(resp.Body.String()), &result)
	require.NoError(t, err)

	require.Equal(t, "APPROVED", result.Status)
	require.Equal(t, "someOID", result.OrderID)
	require.Equal(t, "someMOID", result.MerchantOrderID)
	require.Equal(t, UserMessage, result.UserMessage)
	require.Equal(t, "", result.ErrorMessage)
	require.Equal(t, "", result.PossibleStatuses)
}

// TestPaymentReturn_ErrorMessage tests the data displayed after an unsuccessful deposit that provides and error message
func TestPaymentReturn_ErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{}

	req, err := http.NewRequest(http.MethodGet, "/payment-return?status=ERROR&orderID=someOID&merchantOrderID=someMOID&errorMessage=SomeError", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/payment-return", server.paymentReturn)

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var result PaymentReturn
	err = json.Unmarshal([]byte(resp.Body.String()), &result)
	require.NoError(t, err)

	require.Equal(t, "ERROR", result.Status)
	require.Equal(t, "someOID", result.OrderID)
	require.Equal(t, "someMOID", result.MerchantOrderID)
	require.Equal(t, UserMessage, result.UserMessage)
	require.Equal(t, "SomeError", result.ErrorMessage)
	require.Equal(t, PossibleStatusesLink, result.PossibleStatuses)
}

// TestPaymentReturn_BadRequest tests that if the provided query params are not valid we will throw a 400 Bad Request
func TestPaymentReturn_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{}

	req, err := http.NewRequest(http.MethodGet, "/payment-return?status=ERROR", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/payment-return", server.paymentReturn)

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}
