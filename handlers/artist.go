package handlers

import (
	"sangeet-server/db"
	"sangeet-server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddArtistGroup(app *fiber.App) {
	artistGroup := app.Group("/artists")

	artistGroup.Get("/", getAllArtists)
	artistGroup.Get("/:id", getArtistInfo)
	artistGroup.Post("/", addArtist)
	artistGroup.Patch("/:id", updateArtist)
	artistGroup.Patch("/:id", addTrackToArtist)
	artistGroup.Delete("/:id", deleteArtist)
}

func getAllArtists(c *fiber.Ctx) error {
	coll := db.GetDBCollection("Artists")
	query := c.Query("sort")
	var opts *options.FindOptions
	switch query {
	case "top":
		opts = options.Find().SetSort(bson.D{{"listens", -1}})
	default:
		opts = options.Find().SetSort(bson.D{{}})
	}

	artists := make([]models.Artist, 0)
	cursor, err := coll.Find(c.Context(), bson.M{}, opts)
	if err != nil {
		return c.Status(505).JSON(fiber.Map{"error": err.Error()})
	}

	for cursor.Next(c.Context()) {
		artist := models.Artist{}
		err := cursor.Decode(&artist)
		if err != nil {
			return c.Status(505).JSON(fiber.Map{"error": err.Error()})
		}
		artists = append(artists, artist)
	}

	return c.Status(200).JSON(fiber.Map{"artists": artists})
}

func getArtistInfo(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Artist ID is required"})
	}
	artistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to get artist", "message": err.Error()})
	}

	var artist models.Artist
	coll := db.GetDBCollection("Artists")
	err = coll.FindOne(c.Context(), bson.M{"_id": artistID}).Decode(&artist)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Error finding resource", "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"artist": artist,
	})
}

func addArtist(c *fiber.Ctx) error {
	var artist models.Artist
	if err := c.BodyParser(&artist); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create artist",
			"message": err.Error(),
		})
	}

	coll := db.GetDBCollection("Artists")
	res, err := coll.InsertOne(c.Context(), artist)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create artist",
			"message": err.Error(),
		})
	}

	// return artist id
	return c.Status(201).JSON(fiber.Map{
		"result": res,
	})
}

func updateArtist(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Artist ID required"})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update artist", "message": err.Error()})
	}

	var artist models.Artist
	if err := c.BodyParser(&artist); err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update artist", "message": err.Error()})
	}

	coll := db.GetDBCollection("Artists")
	res, err := coll.UpdateByID(c.Context(), bson.M{"_id": objectID}, artist)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update artist", "message": err.Error()})
	}
	return c.Status(200).
		JSON(fiber.Map{"message": res})
}

/*
ISSUE:
req not idempotent

make idempotent
- [ ] find artist
- [ ] add to array if not added
- [ ] remove if already
*/
func addTrackToArtist(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID is required"})
	}
	albumID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID", "message": err.Error()})
	}
	// verify track id
	var trackIDBody struct {
		ID string `bson:"track_id" json:"track_id"`
	}

	err = c.BodyParser(&trackIDBody)
	if err != nil {
		c.Status(400).
			JSON(fiber.Map{"error": "Something went wrong with your request", "message": err.Error()})
	}
	trackID, err := primitive.ObjectIDFromHex(trackIDBody.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID", "message": err.Error()})
	}

	updateOpts := bson.D{{"$push", bson.D{{"tracks", trackID}}}}
	coll := db.GetDBCollection("Albums")
	res, err := coll.UpdateByID(
		c.Context(),
		bson.M{"_id": albumID},
		updateOpts,
	)

	return c.Status(200).JSON(fiber.Map{"res": res})
}

func deleteArtist(c *fiber.Ctx) error {
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

	coll := db.GetDBCollection("Artists")
	res, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error":   "Could not delete artist",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": res})
}
