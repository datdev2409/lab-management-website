package models

import (
	"time"
)

type Patient struct {
	// ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	ID        string    `json:"id" bson:"_id" db:"id"`
	Name      string    `json:"name" bson:"name" db:"name"`
	YOB       string    `json:"yob" bson:"yob" db:"yob"`
	Gender    string    `json:"gender" bson:"gender" db:"gender"`
	Address   string    `json:"address" bson:"address" db:"address"`
	Phone     string    `json:"phone" bson:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

func NewPatient(name string, yob string, gender string, address string, phone string) *Patient {
	patientId := GenerateRandomID("patient_")
	now := time.Now()
	return &Patient{
		ID:        patientId,
		Name:      name,
		YOB:       yob,
		Gender:    gender,
		Address:   address,
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type PatientQueryOptions struct {
	Keyword string
}

type PatientUpdate struct {
	Name    *string `json:"name,omitempty"`
	YOB     *string `json:"yob,omitempty"`
	Gender  *string `json:"gender,omitempty"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
}

func GetStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
