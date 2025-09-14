package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/auth"
	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Serve the register page
func (h *Handler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.RegisterPage())
}

func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.LoginPage())
}

// RegisterHandler registers a new user
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		AdminToken string `json:"admin_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequestError(AUTH_INVALID_PAYLOAD_ERROR)
	}
	if req.Username == "" || req.Password == "" || req.AdminToken == "" {
		return BadRequestError(AUTH_INVALID_PAYLOAD_ERROR)
	}
	adminToken := os.Getenv("ADMIN_TOKEN")
	if req.AdminToken != adminToken {
		return UnauthorizedError(AUTH_INVALID_ADMIN_TOKEN_ERROR)
	}

	// Check if user already exists
	ctx := r.Context()
	existingUser, err := h.Store.GetUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return BadRequestError(AUTH_USER_EXISTS_ERROR)
	}
	// If error is not user not found, return error
	if err != nil && err.Error() != "user not found" {
		return InternalServerError(err.Error())
	}

	user := models.NewUser(req.Username, req.Password)
	_, err = h.Store.CreateUser(ctx, user)
	if err != nil {
		return InternalServerError(err.Error())
	}

	RespondJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
	return nil
}

// LoginHandler authenticates user and returns JWT token
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromCtx(r.Context())
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("Failed to decode login request", zap.Error(err))
		return BadRequestError(AUTH_INVALID_PAYLOAD_ERROR)
	}
	if req.Username == "" || req.Password == "" {
		return BadRequestError(AUTH_INVALID_PAYLOAD_ERROR)
	}

	ctx := r.Context()
	user, err := h.Store.GetUserByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return UnauthorizedError(AUTH_USER_NOT_FOUND_ERROR)
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if bcryptErr != nil {
		return UnauthorizedError(AUTH_LOGIN_FAILED_ERROR)
	}
	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return InternalServerError("Failed to generate token")
	}

	// Set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(23 * time.Hour), //jwt token valid for 24 hours, to be safe, set cookie to 23 hours
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
	return nil
}
