package handlers

import (
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"lifeline/internal/services"
	"lifeline/pkg/database"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	tokenService services.TokenService
	photoService services.PhotoService
}

func NewAuthHandler(t services.TokenService, p services.PhotoService) *AuthHandler {
	return &AuthHandler{tokenService: t, photoService: p}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	input := new(dto.RegisterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data Format"})
	}
	var count int64
	database.DB.Model(&models.User{}).Where("user_name = ?", strings.ToLower(input.UserName)).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Username is taken"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// Handle Date Parsing (Frontend sends string "YYYY-MM-DD")
	var dob time.Time
	if input.DateOfBirth != "" {
		dob, _ = time.Parse("2006-01-02", input.DateOfBirth)
	} else {
		dob = time.Now() 
	}

	country := input.Country
	if country == "" {
		country = "Morocco"
	}

	//  Map DTO to Model
	user := models.User{
		UserName:     strings.ToLower(input.UserName),
		PasswordHash: string(hashedPassword),
		Name:         input.Name,
		Gender:       input.Gender, 
		BloodGroup:   input.BloodGroup,
		Email:        input.Email,
		PhoneNumber:  input.PhoneNumber,
		DateOfBirth:  dob,
		LastActive:   time.Now(),
		Available:    true, 
		Address: models.Address{
			// City:       input.City,
			Area:       input.Area,
			State:      input.State,
			Country:    country,
			PostalCode: input.PostalCode,
		},
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	database.DB.Model(&user).Association("Roles").Append(&models.Role{Name: "Donor"})

	token, _ := h.tokenService.CreateToken(&user, "")

	return c.JSON(dto.UserResponse{
		UserName: user.UserName,
		Name:     user.Name,
		Token:    token,
	})
}

// Login:  
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	input := new(dto.LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data Format"})
	}

	var user models.User
	err := database.DB.Preload("Photo").Preload("Roles").
		Where("user_name = ? OR email = ?", strings.ToLower(input.UserName), strings.ToLower(input.UserName)).
		First(&user).Error

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"}) 
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	userRole := "Donor"
	for _, r := range user.Roles {
		if r.Name == "Admin" {
			userRole = "Admin"
			break
		}
		if r.Name == "Sponsor" {
			userRole = "Sponsor"
		}
	}

	//  Generate Token
	token, err := h.tokenService.CreateToken(&user, userRole)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	//  Return Response
	return c.JSON(dto.UserResponse{
		UserName: user.UserName,
		Name:     user.Name,
		Gender:   user.Gender,
		Token:    token,
		PhotoURL: user.Photo.URL,
		Role:     userRole, 
	})
}

// GetUserProfile:
func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)

	var user models.User
	database.DB.Preload("Address").Preload("Photo").First(&user, userID)

	response := dto.UserProfile{
		ID:        user.ID,
		Name:      user.Name,
		UserName:  user.UserName,
		// City:      user.Address.City,
		Available: user.Available,
		PhotoURL:  user.Photo.URL,
		BloodGroup:user.BloodGroup,
		// ... Complete other fields
	}

	return c.JSON(response)
}

// Helper:
func getUserIDFromToken(c *fiber.Ctx) uint {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	return uint(claims["nameid"].(float64))
}
