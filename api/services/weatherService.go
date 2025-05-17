package services

import (
	"errors"
	"fmt"
	"time"
	"weatherapi/configs"
	"weatherapi/models"
	"weatherapi/server"
)

type WeatherRecordBody struct {
	RecordedAt  string  `json:"date"`
	Humidity    float32 `json:"humidity"`
	Temperature float32 `json:"temperature"`
}

type rawRecord struct {
	Humidity    float32 `json:"humidity"`
	Temperature float32 `json:"temperature"`
}
type formattedRecord struct {
	Humidity    string `json:"humidity"`
	Temperature string `json:"temperature"`
}
type RecordResponse struct {
	Date      string          `json:"date"`
	Raw       rawRecord       `json:"raw"`
	Formatted formattedRecord `json:"formatted"`
}

func getFormattedRecords(weatherRecords *[]models.Weather, columnsConfig *configs.ColumnsConfig) ([]RecordResponse, error) {
	var dateFormats = map[string]string{
		"YYYY-MM-DD": "2006-01-02",
	}

	dateFormat := dateFormats[columnsConfig.Columns["Date"].Unit]
	if dateFormat == "" {
		return nil, errors.New("invalid date format")
	}

	var results []RecordResponse
	for _, record := range *weatherRecords {
		dateFormatted, _ := time.Parse(time.RFC3339, record.RecordedAt)
		results = append(results, RecordResponse{
			Date: dateFormatted.Format(dateFormat),
			Raw: rawRecord{
				Humidity:    record.Humidity,
				Temperature: record.Temperature,
			},
			Formatted: formattedRecord{
				Humidity:    fmt.Sprintf("%.2f", record.Humidity) + columnsConfig.Columns["Humidity"].Unit,
				Temperature: fmt.Sprintf("%.2f", record.Temperature) + columnsConfig.Columns["Temperature"].Unit,
			},
		})
	}
	return results, nil
}

func GetWeatherRecordsForSingleDay(from string) ([]RecordResponse, error) {
	db := server.GetDb()
	columnsConfig := configs.GetColumns()

	var weatherRecords *[]models.Weather
	db.Where("recorded_at = ?", from).Find(&weatherRecords)

	return getFormattedRecords(weatherRecords, columnsConfig)
}

func GetWeatherRecordsForRange(from string, to string) ([]RecordResponse, error) {
	db := server.GetDb()
	columnsConfig := configs.GetColumns()

	var weatherRecords *[]models.Weather
	db.Where("recorded_at >= ?", from).Where("recorded_at <= ?", to).Find(&weatherRecords)

	return getFormattedRecords(weatherRecords, columnsConfig)
}

func CreateWeatherRecord(record *WeatherRecordBody) (RecordResponse, error) {
	db := server.GetDb()
	columnsConfig := configs.GetColumns()

	weatherRecord := models.Weather{RecordedAt: record.RecordedAt, Humidity: record.Humidity, Temperature: record.Temperature}
	if err := db.Create(&weatherRecord).Error; err != nil {
		return RecordResponse{}, err
	}

	results, err := getFormattedRecords(&[]models.Weather{weatherRecord}, columnsConfig)
	if err != nil {
		return RecordResponse{}, err
	}
	return results[0], nil
}
