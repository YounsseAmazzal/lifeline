package middleware

import (
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"os"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt" 
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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
				"error": "Unauthorized or Token Expired 🚫",
			})
		},
	})
}

func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		userToken := c.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)

		userID := uint(claims["nameid"].(float64))

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
				"error": "Access Denied: You need role " + requiredRole + " ⛔",
			})
		}

		return c.Next()
	}
}

func LogUserActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}

userLocal := c.Locals("user")
if userLocal == nil {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized 🚫",
	})
}

userToken, ok := userLocal.(*jwt.Token)
if !ok {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid token 🚫",
	})
}
		if userToken == nil {
			return nil // ماشي مكونيكطي
		}

		claims := userToken.Claims.(jwt.MapClaims)
		// userID غالباً كيكون float64 فـ JWT map claims
		userIDFloat := claims["nameid"].(float64)
		userID := uint(userIDFloat)

		// 3. تحديث LastActive فالداتابيز
		// كنستعملو GORM UpdateColumn باش نحدثو غير هاد الحقل بلا ما نجبدو اليوزر كامل
		go func(id uint) {
			// درناها وسط Goroutine باش ما نتقلوش على الـ Response ديال اليوزر (Fire and Forget)
			database.DB.Model(&models.User{}).Where("id = ?", id).UpdateColumn("last_active", time.Now())
		}(userID)

		return nil
	}
}
