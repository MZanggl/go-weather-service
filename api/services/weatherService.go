package services

import (
	"fmt"
	"time"
	"weatherapi/configs"
	"weatherapi/models"
	"weatherapi/server"

	"gorm.io/gorm"
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
	var results []RecordResponse
	for _, record := range *weatherRecords {
		dateFormatted, err := time.Parse(columnsConfig.DateFormat, record.RecordedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		results = append(results, RecordResponse{
			Date: dateFormatted.Format(columnsConfig.DateFormat),
			Raw: rawRecord{
				Humidity:    record.Humidity,
				Temperature: record.Temperature,
			},
			Formatted: formattedRecord{
				Humidity:    fmt.Sprintf("%.2f", record.Humidity) + columnsConfig.HumidityFormat,
				Temperature: fmt.Sprintf("%.2f", record.Temperature) + columnsConfig.TemperatureFormat,
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

	var result RecordResponse

	err := db.Transaction(func(tx *gorm.DB) error {
		weatherRecord := models.Weather{
			RecordedAt:  record.RecordedAt,
			Humidity:    record.Humidity,
			Temperature: record.Temperature,
		}
		if err := tx.Create(&weatherRecord).Error; err != nil {
			return fmt.Errorf("error creating record: %v", err)
		}

		results, err := getFormattedRecords(&[]models.Weather{weatherRecord}, columnsConfig)
		if err != nil {
			return fmt.Errorf("error formatting results: %v", err)
		}
		result = results[0]
		return nil
	})

	if err != nil {
		return RecordResponse{}, err
	}
	return result, nil
}
