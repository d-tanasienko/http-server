package config_test

import (
	"httpserver/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort_Default(t *testing.T) {
	port := config.GetPort()
	assert.Equal(t, ":3000", port, "Default port should be :3000")
}

func TestGetPort_EnvironmentVariable(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	port := config.GetPort()
	assert.Equal(t, ":8080", port, "Port from environment variable should be returned")
}

func TestGetBaseUrl_Default(t *testing.T) {
	baseUrl := config.GetBaseUrl()
	assert.Equal(t, "localhost", baseUrl, "Default base URL should be localhost")
}

func TestGetBaseUrl_EnvironmentVariable(t *testing.T) {
	os.Setenv("BASE_URL", "example.com")
	defer os.Unsetenv("BASE_URL")

	baseUrl := config.GetBaseUrl()
	assert.Equal(t, "example.com", baseUrl, "Base URL from environment variable should be returned")
}
