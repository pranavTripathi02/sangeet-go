package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrackArtist struct {
	ID        primitive.ObjectID `bson:"id,omitempty"         json:"id,omitempty"`
	FirstName string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string             `bson:"last_name"            json:"last_name"`
}
type TrackAlbum struct {
	ID    primitive.ObjectID `bson:"id,omitempty"    json:"id,omitempty"`
	Title string             `bson:"title,omitempty" json:"title,omitempty"`
}

type Track struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"       json:"_id,omitempty"`
	Title    string             `bson:"title,omitempty"     json:"title,omitempty"`
	Image    string             `bson:"image,omitempty"     json:"image,omitempty"`
	Artist   TrackArtist        `bson:"artist,omitempty"    json:"artist,omitempty"`
	Album    TrackAlbum         `bson:"album,omitempty"     json:"album,omitempty"`
	Listens  uint16             `bson:"listens,omitempty"   json:"listens,omitempty"`
	TrackNum uint8              `bson:"track_num,omitempty" json:"track_num,omitempty"`
	TrackURI string             `bson:"track_uri,omitempty" json:"track_uri,omitempty"`
}
