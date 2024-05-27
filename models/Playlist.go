package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlaylistTracks struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
}

type Playlist struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"          json:"_id,omitempty"`
	Title        string             `bson:"title,omitempty"        json:"title,omitempty"`
	Image        string             `bson:"image,omitempty"        json:"image,omitempty"`
	Tracks       []PlaylistTracks   `bson:"tracks"                 json:"tracks"`
	Followers    uint8              `bson:"followers"              json:"followers"`
	Listens      uint8              `bson:"listens"                json:"listens"`
	TrendingRank uint8              `bson:"trending_rank"          json:"trending_rank"`
	TracksNum    uint8              `bson:"tracks_numb"            json:"tracks_numb"`
	DurationMS   uint16             `bson:"playlist_duration_ms"   json:"playlist_duration_ms"`
	Link         string             `bson:"playlist_uri,omitempty" json:"playlist_uri,omitempty"`
}
