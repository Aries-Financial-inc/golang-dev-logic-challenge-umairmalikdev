package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/umairmalik/fiber-options-analysis/routes"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}