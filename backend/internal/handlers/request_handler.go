package handlers

import (
	"fmt"
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"lifeline/pkg/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RequestHandler struct{}

func (h *RequestHandler) CreateRequest(c *fiber.Ctx) error {
	// 1. Get User ID from Token
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["nameid"].(float64))

	// 2. Parse Input
	input := new(dto.CreateRequestInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data"})
	}

	// 3. Create Model
	req := models.BloodRequest{
		UserID:       userID,
		BloodType:    input.BloodType,
		IsUrgent:     input.IsUrgent,
		HospitalName: input.HospitalName,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		Status:       models.StatusPending, // Default
		CreatedAt:    time.Now(),
	}

	// 4. Handle Photo Upload (Simple version for now)
	// Hna n9dro n-zido logic dyal photoService mn b3d

	// 5. Save to DB
	if err := database.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create request"})
	}

	return c.JSON(req)
}

func (h *RequestHandler) GetRequests(c *fiber.Ctx) error {
	var requests []models.BloodRequest
	// Jib ghir Approved requests (l-nas ychoufuhom)
	database.DB.Preload("User").Where("status = ?", models.StatusApproved).Find(&requests)
	return c.JSON(requests)
}

//  Get ALL Requests (For Admin)
func (h *RequestHandler) GetAllRequests(c *fiber.Ctx) error {
	var requests []models.BloodRequest
	// Jib l-jdid howa lowel + Ma3loumat User
	database.DB.Preload("User").Order("created_at desc").Find(&requests)
	return c.JSON(requests)
}

//  Approve/Reject Request
func (h *RequestHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	// Body: {"status": "Approved"}
	type UpdateInput struct {
		Status models.RequestStatus `json:"status"`
	}
	input := new(UpdateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var req models.BloodRequest
	if err := database.DB.First(&req, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Request not found"})
	}

	req.Status = input.Status
	database.DB.Save(&req)

	if req.Status == models.StatusApproved {
		go func() {
			// (Mn b3d n-zido l-Geolocalisation hna)
			fmt.Println("Admin Approved! Searching for donors...") // DEBUG 1
			var donors []models.User
			database.DB.Where("blood_group = ?", req.BloodType).Find(&donors)

			    // fmt.Printf("🔍 Found %d donors with type %s\n", len(donors), req.BloodType) // DEBUG 2
			//  Créer Notification l kola wa7ed
			for _, donor := range donors {
				// Ma-nsiftoch l-chakhs li talab (ila kan hwa nit)
				if donor.ID == req.UserID {
				// fmt.Println("   ⏭️ Skipping requester:", donor.UserName) // DEBUG 3
					continue
				}
				notification := models.Notification{
					UserID:    donor.ID,
					Title:     "Urgent Blood Request!",
					Message:   "Someone nearby needs " + req.BloodType + " at " + req.HospitalName,
					Type:      "Urgent",
					IsRead:    false,
					CreatedAt: time.Now(),
				}
				database.DB.Create(&notification)
			}
		}()
	}
	return c.JSON(req)
}

// AcceptRequest: Donor accepts the mission
func (h *RequestHandler) AcceptRequest(c *fiber.Ctx) error {
	//  Get Donor ID
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	donorID := uint(claims["nameid"].(float64))

	reqID := c.Params("id")

	var req models.BloodRequest
	if err := database.DB.First(&req, reqID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Request not found"})
	}

	if req.Status != models.StatusApproved {
		return c.Status(400).JSON(fiber.Map{"error": "Request not available"})
	}

	req.Status = "In_Progress"
	database.DB.Save(&req)

	//  Notify The Requester (L-Mrid)
	go func() {
		var donor models.User
		database.DB.First(&donor, donorID)

		notif := models.Notification{
			UserID:  req.UserID, 
			Title:   "Help is coming! 🚑",
			Message: donor.Name + " has accepted your request and is on the way.",
			Type:    "Success",
			IsRead:  false,
			CreatedAt: time.Now(),
		}
		database.DB.Create(&notif)
	}()

	return c.JSON(fiber.Map{"message": "Thank you! Go save a life."})
}

// Get My Active Request
func (h *RequestHandler) GetMyActiveRequest(c *fiber.Ctx) error {
    // Get User ID
    userToken := c.Locals("user").(*jwt.Token)
    claims := userToken.Claims.(jwt.MapClaims)
    userID := uint(claims["nameid"].(float64))

    var req models.BloodRequest
    // Check if user has a request In_Progress
    err := database.DB.Where("user_id = ? AND status = ?", userID, "In_Progress").First(&req).Error
    
    if err != nil {
        return c.JSON(nil) 
    }
    return c.JSON(req)
}