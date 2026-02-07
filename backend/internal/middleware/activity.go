package middleware

import (
	"fmt"
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"os"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

//  Protected Middleware (Verifies Token)
func Protected() fiber.Handler {
	secret := os.Getenv("TOKEN_KEY")
	if secret == "" {
		secret = "super_secret_default_key_change_me"
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ContextKey: "user", 
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid or Expired Token 🚫",
			})
		},
	})
}

// Role Check Middleware
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: No Token Found"})
		}

		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token Claims"})
		}

		idFloat, ok := claims["nameid"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid User ID in Token"})
		}
		userID := uint(idFloat)

		// Check Role in DB
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		hasRole := false
		for _, role := range user.Roles {
			if strings.EqualFold(role.Name, requiredRole) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Access Denied: You need '%s' role ⛔", requiredRole),
			})
		}

		return c.Next()
	}
}

//  Activity Logger (Update LastActive)
func LogUserActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Execute request first
		err := c.Next()
		if err != nil {
			return err
		}

		// Check if user is logged in
		userLocal := c.Locals("user")
		if userLocal == nil {
			return nil
		}

		userToken, ok := userLocal.(*jwt.Token)
		if !ok {
			return nil 
		}

		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil
		}

		idFloat, ok := claims["nameid"].(float64)
		if !ok {
			return nil
		}
		userID := uint(idFloat)
		go func(id uint) {
			database.DB.Model(&models.User{}).Where("id = ?", id).UpdateColumn("last_active", time.Now())
		}(userID)

		return nil
	}
}