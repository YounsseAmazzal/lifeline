package handlers

import (
	"encoding/json"
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"math"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BankHandler struct{}

func (h *BankHandler) GetBanks(c *fiber.Ctx) error {
	// 1. Pagination Inputs
	page, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	if page < 1 {
		page = 1
	}
	
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize

	// Filters
	bloodGroup := c.Query("bloodGroup")
	bloodGroup, _ = url.QueryUnescape(bloodGroup)
	city := c.Query("city")

	var banks []models.Bank
	var totalCount int64

	// Start Query
	query := database.DB.Model(&models.Bank{})

	// 2. Filters
	if city != "" {
		query = query.Joins("JOIN addresses ON addresses.bank_id = banks.id").
			Where("addresses.city = ?", city)
	}

	if bloodGroup != "" {
		query = query.Joins("JOIN blood_groups ON blood_groups.bank_id = banks.id").
			Where("blood_groups.group = ? AND blood_groups.quantity > 0", bloodGroup)
	}

	// 3. Count (Distinct IDs)
	if err := query.Distinct("banks.id").Count(&totalCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// 4. Check Empty
	if totalCount == 0 {
		return c.JSON([]dto.BankResponse{}) // Return empty array
	}

	// 5. Fetch Data
	err := query.Preload("Address").Preload("BloodGroups").Preload("Photo").
		Order("created_at desc").
		Offset(offset).Limit(pageSize).
		Find(&banks).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch banks"})
	}

	// 6. Pagination Header
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	paginationData, _ := json.Marshal(map[string]interface{}{
		"currentPage":  page,
		"itemsPerPage": pageSize,
		"totalItems":   totalCount,
		"totalPages":   totalPages,
	})
	c.Set("X-Pagination", string(paginationData))
	c.Set("Access-Control-Expose-Headers", "X-Pagination")

	// 7. Mapping
	response := make([]dto.BankResponse, 0)
	
	for _, b := range banks {
		var bgDTOs []dto.BloodGroupDto
		for _, bg := range b.BloodGroups {
			bgDTOs = append(bgDTOs, dto.BloodGroupDto{
				ID: bg.ID, Group: bg.Group, Quantity: bg.Quantity,
			})
		}

		response = append(response, dto.BankResponse{
			ID:          b.ID,
			Name:        b.Name,
			City:        b.Address.City,
			PhoneNumber: b.PhoneNumber,
			LastUpdated: b.CreatedAt,
			PhotoURL:    b.Photo.URL,
			BloodGroups: bgDTOs,
			
			// --- FIX IS HERE ---
			Latitude:    float64(b.Address.Latitude), 
			Longitude:   float64(b.Address.Longitude),
		})
	}

	return c.JSON(response)
}