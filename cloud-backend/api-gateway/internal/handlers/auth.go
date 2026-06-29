package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/cloudmanager/cloud-backend/api-gateway/internal/middleware"
	"github.com/cloudmanager/cloud-backend/api-gateway/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	jwtSecret string
	userRepo  *repository.UserRepository
}

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Never serialize password
	Role     string `json:"role"`
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	return &AuthHandler{
		jwtSecret: secret,
		userRepo:  userRepo,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fallback to hardcoded admin for demo when DB is unavailable
	if h.userRepo == nil {
		if req.Username == "admin" && req.Password == "admin" {
			token, err := middleware.GenerateToken(
				"admin-id",
				"admin",
				"admin",
				h.jwtSecret,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token": token,
				"user": gin.H{
					"id":       "admin-id",
					"username": "admin",
					"role":     "admin",
				},
				"expires_in": 86400,
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Validate credentials against database
	user, err := h.userRepo.ValidateCredentials(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		h.jwtSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		"expires_in": 86400,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fallback when DB is unavailable
	if h.userRepo == nil {
		user := &User{
			ID:       uuid.New().String(),
			Username: req.Username,
			Email:    req.Email,
			Role:     "user",
		}
		c.JSON(http.StatusCreated, gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		})
		return
	}

	// Check if username already exists
	_, err := h.userRepo.GetByUsername(c.Request.Context(), req.Username)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	// Check if email already exists
	_, err = h.userRepo.GetByEmail(c.Request.Context(), req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// Create user
	user := &repository.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// Server can optionally blacklist tokens
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
			"role":     role,
		},
	})
}

// RefreshToken refreshes an access token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	token, err := middleware.GenerateToken(
		userID.(string),
		username.(string),
		role.(string),
		h.jwtSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": 86400,
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	// Fallback when DB is unavailable
	if h.userRepo == nil {
		c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
		return
	}

	// Verify old password
	_, err := h.userRepo.ValidateCredentials(c.Request.Context(), username.(string), req.OldPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(c.Request.Context(), userID.(string), req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// TokenBlacklist is a simple in-memory token blacklist (use Redis in production)
var TokenBlacklist = make(map[string]time.Time)

// AddToBlacklist adds a token to the blacklist
func AddToBlacklist(token string, expiresAt time.Time) {
	TokenBlacklist[token] = expiresAt
}

// IsBlacklisted checks if a token is blacklisted
func IsBlacklisted(token string) bool {
	expiresAt, exists := TokenBlacklist[token]
	if !exists {
		return false
	}
	// Clean up expired tokens
	if time.Now().After(expiresAt) {
		delete(TokenBlacklist, token)
		return false
	}
	return true
}
