package main

import (
	"log"
	"os"
	"sangeet-server/db"
	"sangeet-server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return err
	}
	err := db.Init()
	if err != nil {
		return err
	}

	defer db.Close()
	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	handlers.AddArtistGroup(app)
	handlers.AddTrackGroup(app)
	handlers.AddAlbumGroup(app)
	handlers.AddPlaylistGroup(app)
	handlers.AddAuthGroup(app)

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	app.Listen(":" + port)

	return nil
}
