package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"`
	Role      string    `bson:"role" json:"role"`
	Active    bool      `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewUser(username, password string) *User {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  string(hashPassword),
		Role:      "user", // Default role
		Active:    true,   // Default to active
		CreatedAt: now,
		UpdatedAt: now,
	}
}
