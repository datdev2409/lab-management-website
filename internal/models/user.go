package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"-"`
}

func NewUser(username, password string) *User {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{
		ID:       GenerateRandomID("user_"),
		Username: username,
		Password: string(hashPassword),
	}
}
