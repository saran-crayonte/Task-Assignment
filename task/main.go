package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/saran-crayonte/task/database"
	"github.com/saran-crayonte/task/routes"
)

func main() {
	app := fiber.New()
	database.ConnectDB()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
