package util

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/spf13/viper"
)

// TestLoadConfigFromFile tests loading the configuration from a file.
func TestLoadConfigFromFile(t *testing.T) {
	path := "."

	config, err := LoadConfig(path)

	require.NoError(t, err)

	require.Equal(t, "localhost:8080", config.HTTPServerAddress)
	require.Equal(t, "someSecretKey", config.SecretKey)
	require.Equal(t, "someEndpoint", config.EndpointId)
	require.Equal(t, "someMerchant", config.MerchantId)
}

// TestLoadConfigFromEnv tests loading the configuration from environment variables.
func TestLoadConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("HTTP_SERVER_ADDRESS", "localhost:9090")
	os.Setenv("ZOTA_SECRET_KEY", "anotherSecretKey")
	os.Setenv("ENDPOINT_ID", "anotherEndpoint")
	os.Setenv("MERCHANT_ID", "anotherMerchant")

	// Ensure to clean up environment variables after the test
	defer os.Clearenv()

	viper.Reset()

	config, err := LoadConfig(".")

	require.NoError(t, err)

	require.Equal(t, "localhost:9090", config.HTTPServerAddress)
	require.Equal(t, "anotherSecretKey", config.SecretKey)
	require.Equal(t, "anotherEndpoint", config.EndpointId)
	require.Equal(t, "anotherMerchant", config.MerchantId)
}
