package services

import (
	"fmt"
	"time"
	"weatherapi/configs"
	"weatherapi/models"
	"weatherapi/server"
	"weatherapi/utils"

	"gorm.io/gorm"
)

type WeatherRecordBody struct {
	RecordedAt  string  `json:"date"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

type RawWeatherRecordUnits struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}
type FormattedWeatherRecordUnits struct {
	Humidity    string `json:"humidity"`
	Temperature string `json:"temperature"`
}
type WeatherRecordResponse struct {
	Date      string                      `json:"date"`
	Raw       RawWeatherRecordUnits       `json:"raw"`
	Formatted FormattedWeatherRecordUnits `json:"formatted"`
}

func getFormattedWeatherRecordUnitss(weatherRecords *[]models.Weather, columnsConfig *configs.ColumnsConfig) ([]WeatherRecordResponse, error) {
	var results []WeatherRecordResponse
	for _, record := range *weatherRecords {
		dateFormatted, err := time.Parse(time.RFC3339, record.RecordedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		results = append(results, WeatherRecordResponse{
			Date: dateFormatted.Format(columnsConfig.DateFormat),
			Raw: RawWeatherRecordUnits{
				Humidity:    record.Humidity,
				Temperature: record.Temperature,
			},
			Formatted: FormattedWeatherRecordUnits{
				Humidity:    utils.FormatFloat(record.Humidity, columnsConfig.HumidityFormat),
				Temperature: utils.FormatFloat(record.Temperature, columnsConfig.TemperatureFormat),
			},
		})
	}
	return results, nil
}

func GetWeatherRecordsForSingleDay(from string) ([]WeatherRecordResponse, error) {
	db := server.GetDb()
	columnsConfig := configs.GetColumns()

	var weatherRecords *[]models.Weather
	db.Where("recorded_at = ?", from).Find(&weatherRecords)

	return getFormattedWeatherRecordUnitss(weatherRecords, columnsConfig)
}

func GetWeatherRecordsForRange(from string, to string) ([]WeatherRecordResponse, error) {
	db := server.GetDb()
	columnsConfig := configs.GetColumns()

	var weatherRecords *[]models.Weather
	db.Where("recorded_at >= ?", from).Where("recorded_at <= ?", to).Find(&weatherRecords)

	return getFormattedWeatherRecordUnitss(weatherRecords, columnsConfig)
}

func CreateWeatherRecord(record *WeatherRecordBody) (WeatherRecordResponse, error) {
	db := server.GetDb()
	columnsConfig := configs.GetColumns()

	var result WeatherRecordResponse

	err := db.Transaction(func(tx *gorm.DB) error {
		weatherRecord := models.Weather{
			RecordedAt:  record.RecordedAt,
			Humidity:    record.Humidity,
			Temperature: record.Temperature,
		}
		if err := tx.Create(&weatherRecord).Error; err != nil {
			return fmt.Errorf("error creating record: %v", err)
		}

		results, err := getFormattedWeatherRecordUnitss(&[]models.Weather{weatherRecord}, columnsConfig)
		if err != nil {
			return fmt.Errorf("error formatting results: %v", err)
		}
		result = results[0]
		return nil
	})

	return result, err
}
