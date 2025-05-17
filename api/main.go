package main

import (
	"fmt"
	"log"
	"weatherapi/configs"
	"weatherapi/handlers"

	"weatherapi/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	server.RegisterWebSocket(app)

	app.Get("/weather/:from", handlers.GetWeatherRecordsForSingleDay)
	app.Get("/weather/:from/:to", handlers.GetWeatherRecordsForRange)
	app.Post("/weather", handlers.CreateWeatherRecord)

	conf := configs.Get()
	fmt.Println("Starting server at", conf.AppHost)
	log.Fatal(app.Listen(conf.AppHost))
}
