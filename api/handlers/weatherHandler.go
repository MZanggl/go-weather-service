package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"weatherapi/configs"
	"weatherapi/services"
	"weatherapi/utils"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

var BroadcastFunc = socketio.Broadcast

func GetWeatherRecordsForSingleDay(c *fiber.Ctx) error {
	from := c.Params("from")

	fmt.Println("Received request for singleweather data for single day", from)

	if !utils.IsValidDate(from) {
		log.Println("Invalid 'from' date format:", from)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Request")
	}

	results, err := services.GetWeatherRecordsForSingleDay(from)
	if err != nil {
		log.Println("Error getting weather records:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	return c.Status(fiber.StatusOK).JSON(results)
}

func GetWeatherRecordsForRange(c *fiber.Ctx) error {
	from := c.Params("from")
	to := c.Params("to")

	fmt.Println("Received request for weather data from", from, "to", to)

	if !utils.IsValidDate(from) {
		log.Println("Invalid 'from' date format:", from)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Request")
	}
	if !utils.IsValidDate(to) {
		log.Println("Invalid 'to' date format:", to)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Request")
	}

	results, err := services.GetWeatherRecordsForRange(from, to)
	if err != nil {
		log.Println("Error getting weather records:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	return c.Status(fiber.StatusOK).JSON(results)
}

func CreateWeatherRecord(c *fiber.Ctx) error {
	conf := configs.Get()
	if c.Get("X-Api-Token") != conf.ApiToken {
		log.Println("Invalid or missing API token")
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	record := new(services.WeatherRecordBody)

	if err := c.BodyParser(record); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if !utils.IsValidDate(record.RecordedAt) {
		log.Println("Invalid date format:", record.RecordedAt)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Request")
	}

	fmt.Println("Received record:", record.RecordedAt, record.Humidity, record.Temperature)

	firstRecord, err := services.CreateWeatherRecord(record)
	if err != nil {
		log.Println("Error creating weather record:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	firstRecordJson, err := json.Marshal(firstRecord)
	if err != nil {
		log.Println("Error marshalling record to JSON:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	fmt.Println("Broadcasting record:", string(firstRecordJson))
	BroadcastFunc(firstRecordJson, socketio.TextMessage)

	return c.Status(fiber.StatusCreated).JSON(firstRecord)
}
