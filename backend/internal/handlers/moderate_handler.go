package handlers

import (
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type ModerateHandler struct{}

func (h *ModerateHandler) GetBankModerators(c *fiber.Ctx) error {
	// Simple implementation
	var moderators []models.Moderator
	database.DB.Preload("User").Preload("Bank").Find(&moderators)
	return c.JSON(moderators)
}

// AddModerator: Zid admin l'bank
func (h *ModerateHandler) AddModerator(c *fiber.Ctx) error {
	//... ghadi n9adha 
	return c.SendStatus(200)
}