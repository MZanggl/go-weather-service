package handlers

import (
	"encoding/json"
	"log"
	"weatherapi/configs"
	"weatherapi/services"
	"weatherapi/utils"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

// abstracted to make it mockable in tests
var BroadcastFunc = socketio.Broadcast

func GetWeatherRecordsForSingleDay(c *fiber.Ctx) error {
	from := c.Params("from")

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

	if utils.IsDateInFuture(record.RecordedAt) {
		log.Println("Date is in the future:", record.RecordedAt)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Request")
	}

	log.Println("Received request to create:", record)

	// verify record is not already in the database
	results, err := services.GetWeatherRecordsForSingleDay(record.RecordedAt)
	if err != nil {
		log.Println("Error getting weather records:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	if len(results) > 0 {
		log.Println("Record already exists for date:", record.RecordedAt)
		return c.Status(fiber.StatusConflict).SendString("Record already exists for date")
	}

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

	log.Println("Broadcasting record:", string(firstRecordJson))
	BroadcastFunc(firstRecordJson, socketio.TextMessage)

	return c.Status(fiber.StatusCreated).JSON(firstRecord)
}
