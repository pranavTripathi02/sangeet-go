package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ArtistAlbums struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
}

type ArtistTracks struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
}

type Artist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"        json:"_id,omitempty"`
	FirstName string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string             `bson:"last_name"            json:"last_name"`
	Image     string             `bson:"artist_image"         json:"artist_image"`
	Albums    []ArtistAlbums     `bson:"albums"               json:"albums"`
	Tracks    []ArtistTracks     `bson:"tracks"               json:"tracks"`
	Listens   uint16             `bson:"listens"              json:"listens"`
	URI       string             `bson:"artist_uri,omitempty" json:"artist_uri,omitempty"`
}
