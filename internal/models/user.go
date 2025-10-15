package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id" db:"id"`
	Username  string    `bson:"username" json:"username" db:"username"`
	Password  string    `bson:"password" json:"-" db:"password"`
	Role      string    `bson:"role" json:"role" db:"role"`
	Active    bool      `bson:"active" json:"active" db:"active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at" db:"updated_at"`
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
