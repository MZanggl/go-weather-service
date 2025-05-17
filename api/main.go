package main

import (
	"fmt"
	"log"
	"weatherapi/configs"
	"weatherapi/handlers"

	"weatherapi/server"

	"github.com/gofiber/fiber/v2"
)

func Setup() *fiber.App {
	app := fiber.New()

	server.RegisterWebSocket(app)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong")
	})
	app.Get("/weather/:from", handlers.GetWeatherRecordsForSingleDay)
	app.Get("/weather/:from/:to", handlers.GetWeatherRecordsForRange)
	app.Post("/weather", handlers.CreateWeatherRecord)

	return app
}

func main() {
	app := Setup()

	conf := configs.Get()
	fmt.Println("Starting server at", conf.AppHost)
	log.Fatal(app.Listen(conf.AppHost))
}
