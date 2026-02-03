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
	// Pagination
	page, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	offset := (page - 1) * pageSize

	// Filters
	bloodGroup := c.Query("bloodGroup")
	bloodGroup, _ = url.QueryUnescape(bloodGroup)
	city := c.Query("city")

	var banks []models.Bank
	var totalCount int64

	query := database.DB.Model(&models.Bank{}).Preload("Address").Preload("BloodGroups").Preload("Photo")

	// Filter by Blood Group Availability (Advanced Query)
	if bloodGroup != "" {
		query = query.Joins("JOIN blood_groups ON blood_groups.bank_id = banks.id").
			Where("blood_groups.group = ? AND blood_groups.value > 0", bloodGroup)
	}

	if city != "" {
		query = query.Joins("JOIN addresses ON addresses.bank_id = banks.id").
			Where("addresses.city = ?", city)
	}

	// 💰 BUSINESS LOGIC: Order by Subscription/Activity
	// (هنا فين تقدر تزيد منطق: البنوك اللي مخلصين كيبانو هما اللولين)
	query = query.Order("last_updated desc")

	query.Count(&totalCount)
	query.Offset(offset).Limit(pageSize).Find(&banks)

	// Pagination Header
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	paginationData, _ := json.Marshal(map[string]interface{}{
		"currentPage": page, "itemsPerPage": pageSize, "totalItems": totalCount, "totalPages": totalPages,
	})
	c.Set("Pagination", string(paginationData))
	c.Set("Access-Control-Expose-Headers", "Pagination")

	// Mapping to DTO
	var response []dto.BankResponse
	for _, b := range banks {
		var bgDTOs []dto.BloodGroupDto
		for _, bg := range b.BloodGroups {
			bgDTOs = append(bgDTOs, dto.BloodGroupDto{ID: bg.ID, Group: bg.Group, Quantity: bg.Quantity})
		}

		response = append(response, dto.BankResponse{
			ID: b.ID, Name: b.Name, City: b.Address.City,
			PhoneNumber: b.PhoneNumber, LastUpdated: b.LastUpdated,
			PhotoURL: b.Photo.URL, BloodGroups: bgDTOs,
		})
	}

	return c.JSON(response)
}