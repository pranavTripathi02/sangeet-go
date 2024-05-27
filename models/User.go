package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"             json:"_id,omitempty"`
	FirstName string             `bson:"user_first_name,omitempty" json:"user_first_name,omitempty"`
	LastName  string             `bson:"user_last_name,omitempty"  json:"user_last_name,omitempty"`
	Email     string             `bson:"user_email,omitempty"      json:"user_email,omitempty"`
	Password  string             `bson:"user_password,omitempty"   json:"user_password,omitempty"`
	Role      string             `bson:"user_role,omitempty"       json:"user_role,omitempty"`
}
