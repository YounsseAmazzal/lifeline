package routes

import (
	"lifeline/internal/handlers"
	"lifeline/internal/middleware"
	"lifeline/internal/repository"
	"lifeline/internal/services"
	"lifeline/pkg/database"    
	"github.com/gofiber/fiber/v2"
)



func SetupRoutes(app *fiber.App) {
		// 1. Init Repositories
	userRepo := repository.NewUserRepository(database.DB)

	// 2. Init Services
	tokenService := services.NewTokenService()
	photoService := services.NewPhotoService()

	// 3. Init Handlers (Inject Repos & Services)
	authHandler := handlers.NewAuthHandler(tokenService, photoService)
	userHandler := handlers.NewUserHandler(userRepo) 
	
	bankHandler := &handlers.BankHandler{} 
	adminHandler := &handlers.AdminHandler{}

	api := app.Group("/api")
	api.Post("/account/register", authHandler.Register)
	api.Post("/account/login", authHandler.Login)

	users := api.Group("/users", middleware.Protected())
	users.Get("/", userHandler.GetUsers)       
	users.Get("/:username", userHandler.GetUser)
	
	// Banks
	banks := api.Group("/banks")
	banks.Get("/", bankHandler.GetBanks)

	//profile
	//  routes.go
	api.Get("/account/profile", middleware.Protected(), authHandler.GetUserProfile)
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("Admin"))
	admin.Get("/users-with-roles", adminHandler.GetUsersWithRoles)
	admin.Post("/edit-roles/:username", adminHandler.EditRoles)
}
