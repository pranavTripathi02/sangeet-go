package handlers

import (
	"sangeet-server/db"
	"sangeet-server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTrackGroup(app *fiber.App) {
	trackGroup := app.Group("/tracks")

	trackGroup.Get("/", getAllTracks)
	trackGroup.Get("/:id", getTrackInfo)
	trackGroup.Post("/", addTrack)
	trackGroup.Put("/:id", updateTrack)
	trackGroup.Delete("/:id", deleteTrack)
}

func getAllTracks(c *fiber.Ctx) error {
	coll := db.GetDBCollection("Tracks")

	tracks := make([]models.Track, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(505).JSON(fiber.Map{"error": err.Error()})
	}

	for cursor.Next(c.Context()) {
		track := models.Track{}
		err := cursor.Decode(&track)
		if err != nil {
			return c.Status(505).JSON(fiber.Map{"error": err.Error()})
		}
		tracks = append(tracks, track)
	}

	return c.Status(200).JSON(fiber.Map{"tracks": tracks})
}

func getTrackInfo(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Track ID is required"})
	}
	trackID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update track", "message": err.Error()})
	}

	var track models.Track
	coll := db.GetDBCollection("Tracks")
	err = coll.FindOne(c.Context(), bson.D{{"_id", trackID}}).Decode(&track)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Error finding resource", "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"track": track,
	})
}

// TODO: update artist tracks
// TODO: update album
func addTrack(c *fiber.Ctx) error {
	var track models.Track
	if err := c.BodyParser(&track); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create track",
			"message": err.Error(),
		})
	}

	coll := db.GetDBCollection("Tracks")
	res, err := coll.InsertOne(c.Context(), track)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create track",
			"message": err.Error(),
		})
	}

	// return track id
	return c.Status(201).JSON(fiber.Map{
		"result": res,
	})
}

func updateTrack(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Track ID required"})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update track", "message": err.Error()})
	}

	var track models.Track
	if err := c.BodyParser(&track); err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update track", "message": err.Error()})
	}

	coll := db.GetDBCollection("Tracks")
	res, err := coll.UpdateByID(c.Context(), bson.M{"_id": objectID}, track)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update track", "message": err.Error()})
	}
	return c.Status(200).
		JSON(fiber.Map{"message": res})
}

func deleteTrack(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		c.Status(400).JSON(fiber.Map{
			"error": "ID is required",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
	}

	coll := db.GetDBCollection("Tracks")
	res, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error":   "Could not delete track",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": res})
}
