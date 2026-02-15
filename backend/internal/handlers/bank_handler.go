package handlers

import (
	"encoding/json"
	"fmt"
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type BankHandler struct{}

// GetBanks: Main Handler
func (h *BankHandler) GetBanks(c *fiber.Ctx) error {
	// 1. Pagination Params
	page := getIntQuery(c, "pageNumber", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// 2. Filters Params
	bloodGroup := c.Query("bloodGroup")
	bloodGroup, _ = url.QueryUnescape(bloodGroup)
	city := c.Query("city")

	// Bda query jdida bla ma tkhllet "Address" join daba
	query := database.DB.Model(&models.Bank{})

	if city != "" {
		query = query.Where("LOWER(city) = ?", strings.ToLower(city))
	}

	if bloodGroup != "" {
		query = query.Joins("JOIN blood_groups ON blood_groups.bank_id = banks.id").
			Where("blood_groups.group = ? AND blood_groups.quantity > 0", bloodGroup)
	}

	var totalCount int64
	if err := query.Distinct("banks.id").Count(&totalCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Set Pagination Headers
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	setPaginationHeader(c, page, pageSize, totalCount, totalPages)

	if totalCount == 0 {
		return c.JSON([]dto.BankResponse{})
	}

	var banks []models.Bank
	err := query.
		Select("banks.*"). 
		Group("banks.id"). 
		Preload("Address").
		Preload("BloodGroups").
		Preload("Photo").
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&banks).Error

	if err != nil {
		fmt.Println("Error finding banks:", err) // Debug
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch banks"})
	}

	// --- DEBUG ---
	if len(banks) > 0 {
		fmt.Println(" DEBUG SUCCESS:")
		fmt.Printf("ID: %d, Name: '%s', City: '%s'\n", banks[0].ID, banks[0].Name, banks[0].City)
	}

	response := h.mapToResponse(banks)
	return c.JSON(response)
}


// Helper: Convert Database Models to JSON DTOs
func (h *BankHandler) mapToResponse(banks []models.Bank) []dto.BankResponse {
	response := make([]dto.BankResponse, 0)

	for _, b := range banks {
		var bgDTOs []dto.BloodGroupDto
		for _, bg := range b.BloodGroups {
			bgDTOs = append(bgDTOs, dto.BloodGroupDto{
				ID: bg.ID, Group: bg.Group, Quantity: bg.Quantity,
			})
		}
		fmt.Println("--the name howa hada --",b.Name)
		response = append(response, dto.BankResponse{
			ID:          b.ID,
			Name:        b.Name,
			// City:        b.Address.City,
			PhoneNumber: b.PhoneNumber,
			Email:       b.Email,
			Website:     b.Website,
			LastUpdated: b.CreatedAt, // Or LastUpdated
			PhotoURL:    b.Photo.URL,
			BloodGroups: bgDTOs,
			Latitude:    float64(b.Address.Latitude),
			Longitude:   float64(b.Address.Longitude),
		})
	}
	return response
}

// Helper: Get Int from Query
func getIntQuery(c *fiber.Ctx, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	num, err := strconv.Atoi(val)
	if err != nil || num < 1 {
		return defaultVal
	}
	return num
}

// Helper: Set Headers
func setPaginationHeader(c *fiber.Ctx, page, size int, total int64, pages int) {
	data, _ := json.Marshal(map[string]interface{}{
		"currentPage":  page,
		"itemsPerPage": size,
		"totalItems":   total,
		"totalPages":   pages,
	})
	c.Set("X-Pagination", string(data))
	c.Set("Access-Control-Expose-Headers", "X-Pagination")
}
