package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	password := "testpassword123"

	user := NewUser(username, password)

	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Contains(t, user.ID, "user_")
	assert.Equal(t, username, user.Username)
	assert.NotEmpty(t, user.Password)
	
	// Password should be hashed, not plain text
	assert.NotEqual(t, password, user.Password)
	
	// Verify password hash is valid
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	assert.NoError(t, err, "Password should be correctly hashed")
}

func TestNewUser_PasswordHashing(t *testing.T) {
	username := "testuser"
	password := "mypassword"

	user1 := NewUser(username, password)
	user2 := NewUser(username, password)

	// Same password should generate different hashes (due to salt)
	assert.NotEqual(t, user1.Password, user2.Password, "Same password should generate different hashes")
	
	// Both should be valid for the original password
	err1 := bcrypt.CompareHashAndPassword([]byte(user1.Password), []byte(password))
	err2 := bcrypt.CompareHashAndPassword([]byte(user2.Password), []byte(password))
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}
