package api

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
)

// MockZotaClient is a mock for ZotaClientInterface
type MockZotaClient struct {
	mock.Mock
}

func (m *MockZotaClient) Post(ctx context.Context, endpoint string, body []byte) (*http.Response, error) {
	args := m.Called(ctx, endpoint, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockZotaClient) Get(ctx context.Context, endpoint string, params url.Values) (*http.Response, error) {
	args := m.Called(ctx, endpoint, params)
	return args.Get(0).(*http.Response), args.Error(1)
}
