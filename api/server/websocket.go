package server

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterWebSocket(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		log.Println("WebSocket upgrade request")
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		log.Printf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		log.Printf("Close event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
		log.Printf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	app.Get("/ws/:id", socketio.New(func(kws *socketio.Websocket) {
		userId := kws.Params("id")
		kws.SetAttribute("user_id", userId)
		kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID)), socketio.TextMessage)
	}))
}
