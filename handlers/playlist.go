package handlers

import (
	"net/http"
	"sangeet-server/db"
	"sangeet-server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddPlaylistGroup(app *fiber.App) {
	playlistGroup := app.Group("/playlists")

	playlistGroup.Get("/", getAllPlaylists)
	playlistGroup.Get("/", getAllPlaylists)
	playlistGroup.Get("/:id", getPlaylist)
	playlistGroup.Put("/:id", editPlaylist)
	playlistGroup.Put("/:id/add", addTrackToPlaylist)
	playlistGroup.Post("/", createPlaylist)
}

//	type PlaylistQuery struct {
//		SortBy string
//		Offset uint8
//		Limit  uint8
//	}
func getAllPlaylists(c *fiber.Ctx) error {
	coll := db.GetDBCollection("Playlists")
	var opts *options.FindOptions
	// c.Query("offset")
	var limitResults int64 = 10
	// limitQuery := c.Query("limit")
	//    switch {
	//    case limitQuery>15:
	//        limitResults=20
	//    case limitQuery>25:
	//        limitResults=30
	//    default:
	//        limitResults=10;
	//    }

	sortOrderQuery := c.Query("sort")
	switch sortOrderQuery {
	case "featured":
		opts = options.Find().SetSort(bson.D{{"followers", 1}}).SetLimit(limitResults)
	case "top":
		opts = options.Find().SetSort(bson.D{{"listens", -1}}).SetLimit(limitResults)
	case "trending":
		opts = options.Find().SetSort(bson.D{{"trending_rank", -1}}).SetLimit(limitResults)
	default:
		opts = options.Find().SetSort(bson.D{{}}).SetLimit(limitResults)
	}

	playlists := make([]models.Playlist, 0)
	cursor, err := coll.Find(c.Context(), bson.M{}, opts)
	if err != nil {
		return c.Status(505).JSON(fiber.Map{"error": err.Error()})
	}

	for cursor.Next(c.Context()) {
		playlist := models.Playlist{}
		err := cursor.Decode(&playlist)
		if err != nil {
			return c.Status(505).JSON(fiber.Map{"error": err.Error()})
		}
		playlists = append(playlists, playlist)
	}

	return c.Status(200).JSON(fiber.Map{"playlists": playlists})
}

func getPlaylist(c *fiber.Ctx) error {
	// verify id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing playlist id"})
	}
	playlistId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusNotFound).
			JSON(fiber.Map{"error": "Failed to get playlist", "message": err.Error()})
	}

	var playlist models.Playlist
	coll := db.GetDBCollection("Playlists")
	err = coll.FindOne(c.Context(), bson.D{{"_id", playlistId}}).Decode(&playlist)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Error finding resource", "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"playlist": playlist,
	})
}

func editPlaylist(c *fiber.Ctx) error {
	return nil
}

func addTrackToPlaylist(c *fiber.Ctx) error {
	return nil
}

func createPlaylist(c *fiber.Ctx) error {
	return nil
}
