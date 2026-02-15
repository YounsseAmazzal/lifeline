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

	//geo
	geoHandler := handlers.NewGeoHandler()

	bankHandler := &handlers.BankHandler{}
	adminHandler := &handlers.AdminHandler{}

	//user request bloood
	requestHandler := &handlers.RequestHandler{}

	api := app.Group("/api")
	api.Post("/account/register", authHandler.Register)
	api.Post("/account/login", authHandler.Login)

	users := api.Group("/users", middleware.Protected())
	users.Get("/", userHandler.GetUsers)
	users.Get("/:username", userHandler.GetUser)

	// Banks
	banks := api.Group("/banks")
	banks.Get("/", bankHandler.GetBanks)

	//post and get
	reqs := api.Group("/requests", middleware.Protected())
	reqs.Post("/", requestHandler.CreateRequest) // Create
	reqs.Get("/", requestHandler.GetRequests)
	reqs.Put("/:id/accept", requestHandler.AcceptRequest)
	reqs.Get("/active", requestHandler.GetMyActiveRequest)
	//profile
	//  routes.go
	//cites
	geo := api.Group("/geo")
	geo.Get("/cities", geoHandler.GetCities)
	api.Get("/account/profile", middleware.Protected(), authHandler.GetUserProfile)
	//admin
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("Admin"))
	//requests
	admin.Get("/requests", requestHandler.GetAllRequests)   // Jib kolchi
	admin.Put("/requests/:id", requestHandler.UpdateStatus) // Valider/Refuser

	//notifications
	notifHandler := &handlers.NotificationHandler{}
	api.Get("/notifications", middleware.Protected(), notifHandler.GetMyNotifications)
	api.Put("/notifications/:id/read", middleware.Protected(), notifHandler.MarkAsRead)
	admin.Get("/users-with-roles", adminHandler.GetUsersWithRoles)
	admin.Post("/edit-roles/:username", adminHandler.EditRoles)
}
