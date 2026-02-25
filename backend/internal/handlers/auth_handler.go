package handlers

import (
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"lifeline/internal/services"
	"lifeline/pkg/database"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

	userName := strings.ToLower(c.FormValue("userName"))
	password := c.FormValue("password")
	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phoneNumber")
	blood := c.FormValue("bloodGroup")
	city := c.FormValue("city")
	country := c.FormValue("country")

	// Check username
	var count int64
	database.DB.Model(&models.User{}).
		Where("user_name = ?", userName).
		Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Username is taken"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// ======================
	// HANDLE PHOTO
	// ======================

	var photoPath string

	file, err := c.FormFile("photo")
	if err == nil {

		// Create uploads folder if not exists
		os.MkdirAll("./uploads", os.ModePerm)

		ext := filepath.Ext(file.Filename)
		filename := uuid.New().String() + ext

		photoPath = "/uploads/" + filename

		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save photo"})
		}
	}

	// ======================

	user := models.User{
		UserName:     userName,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Email:        email,
		PhoneNumber:  phone,
		BloodGroup:   blood,
		Photoprofile: photoPath,
		LastActive:   time.Now(),
		Available:    true,
		Address: models.Address{
			City:    city,
			Country: country,
		},
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

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
		ID:         user.ID,
		Name:       user.Name,
		UserName:   user.UserName,
		City:       user.Address.City,
		Available:  user.Available,
		PhotoURL:   user.Photoprofile,
		BloodGroup: user.BloodGroup,
		// ... Complete other fields
	}

	return c.JSON(response)
}

// func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error {

// 	userID := getUserIDFromToken(c)

// 	var user models.User
// 	if err := database.DB.First(&user, userID).Error; err != nil {
// 		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
// 	}

// 	// =========================
// 	// UPDATE TEXT FIELDS
// 	// =========================

// 	user.Name = c.FormValue("name")
// 	user.Email = c.FormValue("email")
// 	user.PhoneNumber = c.FormValue("phone_number")
// 	user.Gender = c.FormValue("gender")
// 	user.BloodGroup = c.FormValue("blood_group")
// 	user.Available = c.FormValue("available") == "true"

// 	// Date parsing
// 	if dobStr := c.FormValue("date_of_birth"); dobStr != "" {
// 		if dob, err := time.Parse("2006-01-02", dobStr); err == nil {
// 			user.DateOfBirth = dob
// 		}
// 	}

// 	// Address
// 	user.Address.City = c.FormValue("city")
// 	user.Address.Area = c.FormValue("area")
// 	user.Address.State = c.FormValue("state")
// 	user.Address.Country = c.FormValue("country")
// 	user.Address.PostalCode = c.FormValue("postal_code")

// 	// =========================
// 	// HANDLE PHOTO UPDATE
// 	// =========================

// 	file, err := c.FormFile("photo")
// 	if err == nil {

// 		// حدف الصورة القديمة إلا كانت
// 		if user.Photoprofile != "" {
// 			oldPath := "." + user.Photoprofile
// 			os.Remove(oldPath)
// 		}

// 		os.MkdirAll("./uploads", os.ModePerm)

// 		ext := filepath.Ext(file.Filename)
// 		filename := uuid.New().String() + ext

// 		newPath := "/uploads/" + filename

// 		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": "Failed to save photo"})
// 		}

// 		user.Photoprofile = newPath
// 	}

// 	// =========================

// 	if err := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error; err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
// 	}

// 	return c.JSON(dto.UserProfile{
// 		ID:         user.ID,
// 		Name:       user.Name,
// 		City:       user.Address.City,
// 		Available:  user.Available,
// 		BloodGroup: user.BloodGroup,
// 		PhotoURL:   user.Photoprofile,
// 	})
// }

// Helper:
func getUserIDFromToken(c *fiber.Ctx) uint {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	return uint(claims["nameid"].(float64))
}
