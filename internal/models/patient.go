package models

type Patient struct {
	ID      string `bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name"`
	YOB     string `json:"yob" bson:"yob"`
	Gender  string `json:"gender" bson:"gender"`
	Address string `json:"address" bson:"address"`
	Phone   string `json:"phone" bson:"phone"`
}
