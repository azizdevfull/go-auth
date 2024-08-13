package controllers

import (
	"go-auth/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ctrl.AuthService.Register(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed, please try again."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": gin.H{
		"username": input.Username,
	}})
}
func (ctrl *AuthController) Profile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	c.JSON(http.StatusOK, gin.H{"data": token})
}

// Login handles user authentication and token generation.
func (ctrl *AuthController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Authenticate user
	authenticated, err := ctrl.AuthService.Authenticate(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := ctrl.AuthService.CreateToken(int(authenticated.ID), authenticated.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with success message and token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"username": authenticated.Username,
			"token":    token,
		},
	})
}
func (ctrl *AuthController) GetUserFromToken(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, err := ctrl.AuthService.FindUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
