package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlbumTracks struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
}

type Album struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"   json:"_id,omitempty"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Image    string             `bson:"image,omitempty" json:"image,omitempty"`
	Tracks   []AlbumTracks      `bson:"tracks"          json:"tracks"`
	Listens  uint16             `bson:"listens"         json:"listens"`
	Link     string             `bson:"link,omitempty"  json:"link,omitempty"`
	Duration uint16             `bson:"album_duration"  json:"album_duration"`
}
