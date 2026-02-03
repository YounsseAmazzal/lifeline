package handlers

import (
	"lifeline/internal/models"
	"github.com/gofiber/fiber/v2"
	"lifeline/pkg/database"
)

// GetImpactReport: 
func GetImpactReport(c *fiber.Ctx) error {
	var totalDonations int64
	database.DB.Model(&models.User{}).Where("last_donation > ?", "2024-01-01").Count(&totalDonations)

	return c.JSON(fiber.Map{
		"campaign": "Ramadan 2024",
		"sponsor": "Maroc Telecom",
		"lives_saved_estimate": totalDonations * 3,
		"active_donors": totalDonations,
	})
}