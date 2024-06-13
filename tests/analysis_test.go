package tests

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"encoding/json"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/gofiber/fiber/v2"
	"github.com/umairmalik/fiber-options-analysis/models"
	"github.com/umairmalik/fiber-options-analysis/routes"
)

func TestOptionsContractModelValidation(t *testing.T) {
	// Example test case to validate the OptionsContract model
	contract := models.OptionsContract{
		Type:          "Call",
		StrikePrice:   100,
		Bid:           10.05,
		Ask:           12.04,
		LongShort:     "long",
		ExpirationDate: time.Now(),
	}

	assert.Equal(t, "Call", contract.Type)
	assert.Equal(t, 100.0, contract.StrikePrice)
	assert.Equal(t, 10.05, contract.Bid)
	assert.Equal(t, 12.04, contract.Ask)
	assert.Equal(t, "long", contract.LongShort)
	assert.False(t, contract.ExpirationDate.IsZero())
}

func TestAnalysisEndpoint(t *testing.T) {
	// Set up Fiber app and define routes
	app := fiber.New()
	routes.SetupRoutes(app)

	// Create a request to the /analyze endpoint with a sample contract
	requestPayload := `[{
		"type": "Call",
		"strike_price": 100,
		"bid": 10.05,
		"ask": 12.04,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}]`
	req := httptest.NewRequest("POST", "/analyze", strings.NewReader(requestPayload))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request and check for no errors
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response and check the presence of expected keys
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "xy_values")
	assert.Contains(t, response, "max_profit")
	assert.Contains(t, response, "max_loss")
	assert.Contains(t, response, "break_even_points")
}

func TestIntegration(t *testing.T) {
	// Integration test for the /analyze endpoint
	app := fiber.New()
	routes.SetupRoutes(app)

	// Sample request with a call option contract
	requestPayload := `[{
		"type": "Call",
		"strike_price": 100,
		"bid": 10.05,
		"ask": 12.04,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}]`
	req := httptest.NewRequest("POST", "/analyze", strings.NewReader(requestPayload))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request and verify the response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse and validate the JSON response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "xy_values")
	assert.Contains(t, response, "max_profit")
	assert.Contains(t, response, "max_loss")
	assert.Contains(t, response, "break_even_points")
}

func TestMultipleContracts(t *testing.T) {
	// Test handling multiple option contracts
	app := fiber.New()
	routes.SetupRoutes(app)

	// Sample request with multiple contracts
	requestPayload := `[{
		"type": "Call",
		"strike_price": 100,
		"bid": 10.05,
		"ask": 12.04,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}, {
		"type": "Put",
		"strike_price": 105,
		"bid": 16,
		"ask": 18,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}]`
	req := httptest.NewRequest("POST", "/analyze", strings.NewReader(requestPayload))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request and validate the response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response and check for expected keys
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "xy_values")
	assert.Contains(t, response, "max_profit")
	assert.Contains(t, response, "max_loss")
	assert.Contains(t, response, "break_even_points")
}

func TestInvalidData(t *testing.T) {
	// Test with invalid data in the request
	app := fiber.New()
	routes.SetupRoutes(app)

	// Request with an invalid option type
	requestPayload := `[{
		"type": "InvalidType",
		"strike_price": 100,
		"bid": 10.05,
		"ask": 12.04,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}]`
	req := httptest.NewRequest("POST", "/analyze", strings.NewReader(requestPayload))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request and check for a bad request status
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Decode the error response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestMixedOptionTypes(t *testing.T) {
	// Test handling a mix of call and put option contracts
	app := fiber.New()
	routes.SetupRoutes(app)

	// Sample request with mixed option types
	requestPayload := `[{
		"type": "Call",
		"strike_price": 100,
		"bid": 10.05,
		"ask": 12.04,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}, {
		"type": "Put",
		"strike_price": 105,
		"bid": 16,
		"ask": 18,
		"long_short": "short",
		"expiration_date": "2025-12-17T00:00:00Z"
	}, {
		"type": "Call",
		"strike_price": 110,
		"bid": 20.10,
		"ask": 22.05,
		"long_short": "long",
		"expiration_date": "2025-12-17T00:00:00Z"
	}, {
		"type": "Put",
		"strike_price": 95,
		"bid": 14,
		"ask": 15.50,
		"long_short": "short",
		"expiration_date": "2025-12-17T00:00:00Z"
	}]`
	req := httptest.NewRequest("POST", "/analyze", strings.NewReader(requestPayload))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request and check for success
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response and verify expected data
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "xy_values")
	assert.Contains(t, response, "max_profit")
	assert.Contains(t, response, "max_loss")
	assert.Contains(t, response, "break_even_points")
}