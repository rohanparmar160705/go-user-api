/*
Package routes defines the API endpoints and maps them to handler functions.
It uses GoFiber's routing mechanism to group routes and apply middleware if needed.
*/
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rohanparmar/go-user-api/internal/handler"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users", userHandler.ListUsers)
	app.Get("/users/:id", userHandler.GetUser)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}

