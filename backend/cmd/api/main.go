package main

import (
	"lifeline/internal/middleware"
	"lifeline/internal/models"
	"lifeline/internal/routes"
	"lifeline/pkg/database"
	"log"
	"os"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	fmt.Println("localhost:8080")
	// Database Connection & Seed
	database.Connect()
	database.DB.AutoMigrate(&models.Region{}, &models.City{})
	// SeedMoroccanCities
	database.SeedMoroccanCities(database.DB) 
	//seedamiandRoles
	database.SeedRolesAndAdmin(database.DB) 
	// banks
	database.SeedBanks(database.DB)          
	//fake users 
	database.SeedFakeUsers(database.DB)
	//  Fiber Config (Startup.cs ConfigureServices)
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middlewares (Startup.cs Configure)
	app.Use(logger.New()) // Logging
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*", 
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Pagination", 
	}))
	// Custom Activity Logger
	app.Use(middleware.LogUserActivity())

	//  Routes
	routes.SetupRoutes(app)
	app.Static("/uploads", "./uploads")	//  Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
