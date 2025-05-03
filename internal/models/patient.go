package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Patient struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	YOB       string        `json:"yob" bson:"yob"`
	Gender    string        `json:"gender" bson:"gender"`
	Address   string        `json:"address" bson:"address"`
	Phone     string        `json:"phone" bson:"phone"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type PatientQueryOptions struct {
	Keyword string
}
