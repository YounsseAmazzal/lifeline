package handlers

import (
	"encoding/json"
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"lifeline/internal/repository" 

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" 
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(r *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

// GetUsers:
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	params := dto.UserParams{}
	if err := c.QueryParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid params"})
	}

	params.SetDefaults()

	// Helper function call
	params.CurrentUserName = getCurrentUserName(c)

	// Call Repository via h.repo
	users, pagination, err := h.repo.GetUsers(params)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error: " + err.Error()})
	}

	// Pagination Header
	paginationJSON, _ := json.Marshal(pagination)
	c.Set("Pagination", string(paginationJSON))
	c.Set("Access-Control-Expose-Headers", "Pagination")

	// Mapping
	var response []dto.MemberResponse
	for _, u := range users {
		response = append(response, mapUserToMemberDto(u))
	}

	return c.JSON(response)
}

// GetUser:
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.repo.GetUser(username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(mapUserToMemberDto(*user))
}

// --- Helpers ---

func getCurrentUserName(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if name, ok := claims["unique_name"].(string); ok {
		return name
	}
	return ""
}

func mapUserToMemberDto(user models.User) dto.MemberResponse {
	photoUrl := ""
	if user.Photo.URL != "" {
		photoUrl = user.Photo.URL
	}

	city := "Unknown"
	if user.Address.City != "" {
		city = user.Address.City
	}

	return dto.MemberResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Name:     user.Name,
		PhotoURL: photoUrl,
		Gender:   user.Gender,
		City:       city,
		BloodGroup: user.BloodGroup,
		LastActive: user.LastActive.Format("2006-01-02"),
		Available:  user.Available,
	}
}
