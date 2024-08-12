package main

import (
	"go-auth/models"
	"go-auth/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "user=postgres password=postgres dbname=postgres host=localhost port=5440 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	r := gin.Default()
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database")
	}
	routes.SetupRoutes(r, db)

	r.Run(":8080")
}
