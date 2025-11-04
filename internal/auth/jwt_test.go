package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	// Set up JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	userID := "user_123"

	token, err := GenerateJWT(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	// Set up JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("valid token", func(t *testing.T) {
		userID := "user_456"
		token, err := GenerateJWT(userID)
		require.NoError(t, err)

		validatedUserID, err := ValidateJWT(token)
		require.NoError(t, err)
		assert.Equal(t, userID, validatedUserID)
	})

	t.Run("invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		_, err := ValidateJWT(invalidToken)
		assert.Error(t, err)
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := ValidateJWT("")
		assert.Error(t, err)
	})

	t.Run("token with wrong secret", func(t *testing.T) {
		// Generate token with one secret
		userID := "user_789"
		token, err := GenerateJWT(userID)
		require.NoError(t, err)

		// Try to validate with different secret
		os.Setenv("JWT_SECRET", "different-secret")
		defer os.Setenv("JWT_SECRET", "test-secret-key-for-testing")

		_, err = ValidateJWT(token)
		assert.Error(t, err)
	})

	t.Run("malformed token", func(t *testing.T) {
		malformedToken := "not-a-jwt-token"
		_, err := ValidateJWT(malformedToken)
		assert.Error(t, err)
	})
}

func TestJWTRoundTrip(t *testing.T) {
	// Set up JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	testCases := []string{
		"user_1",
		"user_2",
		"admin_123",
		"test_user",
	}

	for _, userID := range testCases {
		t.Run(userID, func(t *testing.T) {
			token, err := GenerateJWT(userID)
			require.NoError(t, err)

			validatedUserID, err := ValidateJWT(token)
			require.NoError(t, err)
			assert.Equal(t, userID, validatedUserID)
		})
	}
}
