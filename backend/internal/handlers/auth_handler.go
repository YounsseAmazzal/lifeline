package handlers

import (
	"fmt"
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

// Register: تسجيل مستخدم جديد
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// 1. Parsing Input
	input := new(dto.RegisterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data"})
	}

	// 2. Validation (Check if username exists)
	var count int64
	database.DB.Model(&models.User{}).Where("user_name = ?", strings.ToLower(input.UserName)).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Username is taken"})
	}

	// 3. Hash Password (بديل UserManager.CreateAsync)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 4. Map DTO to Model
	user := models.User{
		UserName:     strings.ToLower(input.UserName),
		PasswordHash: string(hashedPassword),
		Name:         input.Name,
		Gender:       input.Gender,
		BloodGroup:   input.BloodGroup,
		Email:        input.Email,
		PhoneNumber:  input.PhoneNumber,
		LastActive:   time.Now(),
		Available:    true,
		Address: models.Address{ // Create Address Automatically
			City:       input.City,
			Area:       input.Area,
			State:      input.State,
			Country:    input.Country,
			PostalCode: input.PostalCode,
		},
		Photo: models.Photo{}, // Empty Photo initially
	}

	// 5. Save to DB
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	// 6. Assign Role "Member"
	// (خاصنا نكونو ديجا عمرنا جدول Roles، غانفترضوه كاين)
	// database.DB.Model(&user).Association("Roles").Append(&models.Role{Name: "Member"})

	// 7. Generate Token
	token, _ := h.tokenService.CreateToken(&user, "")

	// 8. Return UserDto
	return c.JSON(dto.UserResponse{
		UserName: user.UserName,
		Name:     user.Name,
		Gender:   user.Gender,
		Token:    token,
		PhotoURL: "", // باقي مادارش تصويرة
	})
}

// Login: تسجيل الدخول
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	input := new(dto.LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	var user models.User
	// نجيبو اليوزر مع التصويرة والرولز (Preload)
	result := database.DB.Preload("Photo").Preload("Roles").
		Where("user_name = ?", strings.ToLower(input.UserName)).First(&user)

	if result.Error != nil {
		fmt.Print("dsds")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user"})
	}

	// تحقق من الباسورد (بديل CheckPasswordSignInAsync)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		fmt.Print("Hash")
		return c.Status(401).JSON(fiber.Map{"error": "Wrong password"})
	}
	
	token, _ := h.tokenService.CreateToken(&user, "")
	
	return c.JSON(dto.UserResponse{
		UserName: user.UserName,
		Name:     user.Name,
		Gender:   user.Gender,
		Token:    token,
		PhotoURL: user.Photo.URL,
	})
}

// GetUserProfile: (Authorize required)
func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)

	var user models.User
	database.DB.Preload("Address").Preload("Photo").First(&user, userID)

	// Map to UserProfileDto (Manual Mapping for simplicity)
	response := dto.UserProfile{
		ID:          user.ID,
		Name:        user.Name,
		UserName:    user.UserName,
		City:        user.Address.City,
		Available:   user.Available,
		PhotoURL:    user.Photo.URL,
		// ... Complete other fields
	}

	return c.JSON(response)
}

// Helper: استخراج UserID من التوكن (Utility)
func getUserIDFromToken(c *fiber.Ctx) uint {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	return uint(claims["nameid"].(float64))
}