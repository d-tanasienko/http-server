package config

import "os"

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return ":" + port
}

func GetBaseUrl() string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "localhost"
	}

	return baseUrl
}
