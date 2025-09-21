package models

import (
	"time"
)

type Doctor struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Phone     string    `json:"phone" bson:"phone"`
	Address   string    `json:"address" bson:"address"` // Optional field
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func NewDoctor(name string, phone string, address string) *Doctor {
	doctorId := GenerateRandomID("doctor_")
	now := time.Now()
	return &Doctor{
		ID:        doctorId,
		Name:      name,
		Phone:     phone,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type DoctorQueryOptions struct {
	Keyword string
}

type DoctorUpdate struct {
	Name    *string `json:"name,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
}

type CreateDoctorRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address,omitempty"` // Optional
}

type UpdateDoctorRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address,omitempty"` // Optional
}

func (d Doctor) GetID() string {
	return d.ID
}