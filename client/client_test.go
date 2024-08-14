package client

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestZotaClient_Get tests the GET request functionality.
func TestZotaClient_Get(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/someEndpoint", r.URL.Path)
		require.Equal(t, "param1=value1&param2=value2", r.URL.RawQuery)
		require.Equal(t, "someSecretKey", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer mockServer.Close()

	client := &ZotaClient{
		client:  &http.Client{},
		headers: map[string]string{"Authorization": "someSecretKey"},
		baseUrl: mockServer.URL,
	}

	params := url.Values{}
	params.Add("param1", "value1")
	params.Add("param2", "value2")

	resp, err := client.Get(context.Background(), "someEndpoint", params)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestZotaClient_Post tests the POST request functionality.
func TestZotaClient_Post(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/someEndpoint", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, "someSecretKey", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer mockServer.Close()

	client := &ZotaClient{
		client: &http.Client{},
		headers: map[string]string{
			"Authorization": "someSecretKey",
			"Content-Type":  "application/json",
		},
		baseUrl: mockServer.URL,
	}

	body := []byte(`{"name": "test"}`)
	resp, err := client.Post(context.Background(), "someEndpoint", body)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

// TestZotaClient_Non2xxResponse tests the handling of non-2xx responses.
func TestZotaClient_Non2xxResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer mockServer.Close()

	client := &ZotaClient{
		client:  &http.Client{},
		headers: map[string]string{"Authorization": "someSecretKey"},
		baseUrl: mockServer.URL,
	}

	body := []byte(`{"name": "test"}`)
	resp, err := client.Post(context.Background(), "someEndpoint", body)
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Contains(t, err.Error(), "received non-2xx status code")
}
