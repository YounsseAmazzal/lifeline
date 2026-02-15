package handlers

import (
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type NotificationHandler struct{}

// GET /api/notifications
func (h *NotificationHandler) GetMyNotifications(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["nameid"].(float64))

	var notifs []models.Notification
	// Jib ghir li ma-mqryinch (IsRead = false) aw akher 10
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&notifs)

	return c.JSON(notifs)
}

// PUT /api/notifications/:id/read (Mark as Read)
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true)
	return c.SendStatus(200)
}