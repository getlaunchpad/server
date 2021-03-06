package tests

import (
	"net/http"
	"testing"
)

func TestMainRouteStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestHealthStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/status/health", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestReadinessStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/status/readiness", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
