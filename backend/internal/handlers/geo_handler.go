package handlers

import (
	"lifeline/internal/models"
	"lifeline/pkg/database"

	"github.com/gofiber/fiber/v2"
)

type GeoHandler struct{}

func NewGeoHandler() *GeoHandler {
	return &GeoHandler{}
}

// GetCities:
func (h *GeoHandler) GetCities(c *fiber.Ctx) error {
	var cities []models.City

	result := database.DB.Order("ville ASC").Find(&cities)
	// log.Fatal(result)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch cities"})
	}

	return c.JSON(cities)
}
