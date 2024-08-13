package routes

import (
	"go-auth/controllers"
	middleware "go-auth/middlewares"
	"go-auth/repositories"
	"go-auth/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}

	profileRoutes := r.Group("/profile")

	profileRoutes.Use(middleware.JWTAuthMiddleware([]byte("secret-key")))
	{
		profileRoutes.GET("", authController.GetUserFromToken)
	}
	// Apply middleware to routes if needed
	// r.Use(middlewares.AuthMiddleware())
}
