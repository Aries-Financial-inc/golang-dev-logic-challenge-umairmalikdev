package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/umairmalik/fiber-options-analysis/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/analyze", controllers.AnalyzeOptionsContracts)
}