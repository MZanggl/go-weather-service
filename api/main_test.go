package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"weatherapi/handlers"
	"weatherapi/models"
	"weatherapi/server"
	"weatherapi/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPingRoute(t *testing.T) {
	app := Setup()
	req, _ := http.NewRequest("GET", "/ping", nil)
	res, err := app.Test(req, -1)

	// Validate response
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Pong", string(body))
}

func prepareTestDB() *gorm.DB {
	db := server.GetDb()
	db.Migrator().DropTable(&models.Weather{})
	db.AutoMigrate(&models.Weather{})
	return db
}

func TestCreateWeatherRoute(t *testing.T) {
	app := Setup()

	t.Run("weather creation endpoint requires a token", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/weather", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 401, res.StatusCode)
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, "Unauthorized", string(body))
	})

	t.Run("weather creation endpoint fails when passing an invalid date", func(t *testing.T) {
		prepareTestDB()

		requestBody := `{"date":"invalid date!","humidity":60.98765,"temperature":25.98765}`
		req, _ := http.NewRequest("POST", "/weather", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Api-Token", "abcdef")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, "Invalid Request", string(body))
	})

	t.Run("weather creation endpoint fails when passing insufficient inputs", func(t *testing.T) {
		prepareTestDB()

		requestBody := `{"temperature":25.98765}`
		req, _ := http.NewRequest("POST", "/weather", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Api-Token", "abcdef")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, "Invalid Request", string(body))
	})

	t.Run("weather creation endpoint saves data to db, returns formatted record, and broadcasts a websocket message", func(t *testing.T) {
		db := prepareTestDB()

		// mock socketio.Broadcast
		original := handlers.BroadcastFunc
		websocketEvent := []byte{}
		handlers.BroadcastFunc = func(event []byte, mType ...int) {
			websocketEvent = event
		}
		defer func() { handlers.BroadcastFunc = original }()

		requestBody := `{"date":"2024-06-01","humidity":60.98765,"temperature":25.98765}`
		req, _ := http.NewRequest("POST", "/weather", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Api-Token", "abcdef")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 201, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		var actual services.RecordResponse
		err = json.Unmarshal(body, &actual)
		assert.Nil(t, err)

		expected := services.RecordResponse{Date: "2024-06-01", Raw: services.RawWeatherRecordUnits{Humidity: 60.98765, Temperature: 25.98765}, Formatted: services.FormattedWeatherRecordUnits{Humidity: "60.99%", Temperature: "25.99°C"}}
		assert.Equal(t, expected, actual)

		// Validate record exists in the database
		var weatherRecords []models.Weather
		err = db.Where("recorded_at = ?", "2024-06-01").Find(&weatherRecords).Error
		assert.Nil(t, err)
		assert.Equal(t, 1, len(weatherRecords))
		assert.Equal(t, 60.98765, weatherRecords[0].Humidity)
		assert.Equal(t, 25.98765, weatherRecords[0].Temperature)
		assert.Equal(t, "2024-06-01", weatherRecords[0].RecordedAt)

		// validate websocket message
		var actualWebsocketEvent services.RecordResponse
		err = json.Unmarshal([]byte(websocketEvent), &actualWebsocketEvent)
		assert.Nil(t, err)
		assert.Equal(t, expected, actualWebsocketEvent)
	})
}

func TestGetWeatherRecordForSingleDayRoute(t *testing.T) {
	app := Setup()

	t.Run("fails when passing an invalid date", func(t *testing.T) {
		prepareTestDB()

		req, _ := http.NewRequest("GET", "/weather/2025-0101T00:00:00", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, "Invalid Request", string(body))
	})

	t.Run("returns all records for the given day", func(t *testing.T) {
		db := prepareTestDB()

		db.Create(&models.Weather{RecordedAt: "2025-01-01", Humidity: 60.98765, Temperature: 25.98765})
		db.Create(&models.Weather{RecordedAt: "2025-01-02", Humidity: 60.98765, Temperature: 25.98765})

		req, _ := http.NewRequest("GET", "/weather/2025-01-01", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 200, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		var actual []services.RecordResponse
		err = json.Unmarshal(body, &actual)
		assert.Nil(t, err)

		expected := []services.RecordResponse{
			{Date: "2025-01-01", Raw: services.RawWeatherRecordUnits{Humidity: 60.98765, Temperature: 25.98765}, Formatted: services.FormattedWeatherRecordUnits{Humidity: "60.99%", Temperature: "25.99°C"}},
		}
		assert.Equal(t, expected, actual)
	})
}

func TestGetWeatherRecordForRangeRoute(t *testing.T) {
	app := Setup()

	t.Run("fails when passing an invalid date", func(t *testing.T) {
		prepareTestDB()

		req, _ := http.NewRequest("GET", "/weather/2025-01-01T00:00:00/2025-01-03T00:00:00", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, "Invalid Request", string(body))
	})

	t.Run("returns all records for the given day", func(t *testing.T) {
		db := prepareTestDB()

		db.Create(&models.Weather{RecordedAt: "2025-01-01", Humidity: 60.98765, Temperature: 25.98765})
		db.Create(&models.Weather{RecordedAt: "2025-01-02", Humidity: 60.98765, Temperature: 25.98765})
		db.Create(&models.Weather{RecordedAt: "2025-01-03", Humidity: 60.98765, Temperature: 25.98765})
		// won't be included
		db.Create(&models.Weather{RecordedAt: "2025-01-04", Humidity: 60.98765, Temperature: 25.98765})

		req, _ := http.NewRequest("GET", "/weather/2025-01-01/2025-01-03", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		// Validate response
		assert.Nil(t, err)
		assert.Equal(t, 200, res.StatusCode)
		body, _ := io.ReadAll(res.Body)

		var actual []services.RecordResponse
		err = json.Unmarshal(body, &actual)
		assert.Nil(t, err)

		expected := []services.RecordResponse{
			{Date: "2025-01-01", Raw: services.RawWeatherRecordUnits{Humidity: 60.98765, Temperature: 25.98765}, Formatted: services.FormattedWeatherRecordUnits{Humidity: "60.99%", Temperature: "25.99°C"}},
			{Date: "2025-01-02", Raw: services.RawWeatherRecordUnits{Humidity: 60.98765, Temperature: 25.98765}, Formatted: services.FormattedWeatherRecordUnits{Humidity: "60.99%", Temperature: "25.99°C"}},
			{Date: "2025-01-03", Raw: services.RawWeatherRecordUnits{Humidity: 60.98765, Temperature: 25.98765}, Formatted: services.FormattedWeatherRecordUnits{Humidity: "60.99%", Temperature: "25.99°C"}},
		}
		assert.Equal(t, expected, actual)
	})
}
