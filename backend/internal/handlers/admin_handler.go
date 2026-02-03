package handlers

import (
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct{}

// GetUsersWithRoles: Admin kichouf kulchi m3a roles
func (h *AdminHandler) GetUsersWithRoles(c *fiber.Ctx) error {
	var users []models.User
	// Preload Roles
	database.DB.Preload("Roles").Find(&users)
	
	var response []fiber.Map
	for _, u := range users {
		var roles []string
		for _, r := range u.Roles {
			roles = append(roles, r.Name)
		}
		response = append(response, fiber.Map{
			"id": u.ID,
			"username": u.UserName,
			"roles": roles,
		})
	}
	return c.JSON(response)
}

// EditRoles: Admin kibdel role dyal user
func (h *AdminHandler) EditRoles(c *fiber.Ctx) error {
	username := c.Params("username")
	var rolesInput struct {
		Roles []string `json:"roles"`
	}
	if err := c.BodyParser(&rolesInput); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data"})
	}

	var user models.User
	if err := database.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	database.DB.Model(&user).Association("Roles").Clear()
	
	var newRoles []models.Role
	database.DB.Where("name IN ?", rolesInput.Roles).Find(&newRoles)
	
	database.DB.Model(&user).Association("Roles").Append(newRoles)

	return c.JSON(rolesInput.Roles)
}