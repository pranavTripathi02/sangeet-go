package handlers

import (
	"sangeet-server/db"
	"sangeet-server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAlbumGroup(app *fiber.App) {
	albumGroup := app.Group("/albums")

	albumGroup.Get("/", getAllAlbums)
	albumGroup.Get("/:id", getAlbumInfo)
	albumGroup.Patch("/:id", editAlbumInfo)
	albumGroup.Patch("/:id/add", addTrackToAlbum)
	albumGroup.Post("/", createAlbum)
}

func getAllAlbums(c *fiber.Ctx) error {
	coll := db.GetDBCollection("Albums")

	albums := make([]models.Album, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(505).JSON(fiber.Map{"error": err.Error()})
	}

	for cursor.Next(c.Context()) {
		album := models.Album{}
		err := cursor.Decode(&album)
		if err != nil {
			return c.Status(505).JSON(fiber.Map{"error": err.Error()})
		}
		albums = append(albums, album)
	}

	return c.Status(200).JSON(fiber.Map{"albums": albums})
}

func getAlbumInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing album id"})
	}
	albumID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid album ID", "message": err.Error()})
	}

	var album models.Album
	coll := db.GetDBCollection("Albums")
	err = coll.FindOne(c.Context(), bson.D{{"_id", albumID}}).Decode(&album)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Could not fetch album", "message": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"album": album})
}

func editAlbumInfo(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID is required"})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID", "message": err.Error()})
	}
	// verify update body
	var updateOpts models.Album
	err = c.BodyParser(&updateOpts)
	if err != nil {
		c.Status(400).
			JSON(fiber.Map{"error": "Something went wrong with your request", "message": err.Error()})
	}

	coll := db.GetDBCollection("Albums")
	res, err := coll.UpdateByID(c.Context(), bson.M{"_id": objectID}, updateOpts)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Could not fetch album", "message": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"res": res})
}

/*
ISSUE:
req not idempotent

make idempotent
- [ ] find album
- [ ] add to array if not added
- [ ] remove if already
*/
type AddTrackToAlbumDTO struct {
	ID string `bson:"_id" json:"_id"`
}

func addTrackToAlbum(c *fiber.Ctx) error {
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
	var trackIDBody AddTrackToAlbumDTO
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

func createAlbum(c *fiber.Ctx) error {
	var album models.Album
	if err := c.BodyParser(&album); err != nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "Failed to create album", "message": err.Error()})
	}

	coll := db.GetDBCollection("Albums")
	res, err := coll.InsertOne(c.Context(), album)
	if err != nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "Failed to create album", "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"res": res})
}
