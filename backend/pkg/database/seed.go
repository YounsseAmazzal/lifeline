package database

import (
	"encoding/json"
	"fmt"
	"lifeline/internal/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- Structs for JSON Parsing ---
type JsonRegion struct {
	ID     string `json:"id"`
	Region string `json:"region"`
}

type JsonCity struct {
	ID       string `json:"id"`
	Ville    string `json:"ville"`
	RegionID string `json:"region"`
}

type JsonAddress struct {
	Area       string  `json:"Area"`
	City      string    `json:"city"`
	State      string  `json:"State"`
	PostalCode string  `json:"PostalCode"`
	Latitude   float64 `json:"Latitude"`
	Longitude  float64 `json:"Longitude"`
}

type JsonBloodGroup struct {
	Group string `json:"Group"`
	Value int    `json:"Value"`
}

type JsonBank struct {
	Name        string           `json:"Name"`
	City        string          `json:"City"`
	PhoneNumber string           `json:"PhoneNumber"`
	Email       string           `json:"Email"`
	Website     string           `json:"Website"`
	Address     JsonAddress      `json:"Address"`
	BloodGroups []JsonBloodGroup `json:"BloodGroups"`
}

// --- Seeding Functions ---

func SeedMoroccanCities(db *gorm.DB) {
	var count int64
	db.Model(&models.Region{}).Count(&count)
	if count > 0 {
		fmt.Println("🇲🇦 Moroccan Data already exists. Skipping.")
		return
	}

	fmt.Println("⏳ Seeding Regions & Cities...")

	dataRegion, err := os.ReadFile("assets/sql-moroccan-cities/json/region.json")
	if err != nil {
		fmt.Println("❌ Error reading region.json:", err)
		return
	}

	var jsonRegions []JsonRegion
	json.Unmarshal(dataRegion, &jsonRegions)

	regionMap := make(map[string]uint)

	for _, r := range jsonRegions {
		region := models.Region{Region: r.Region}
		if err := db.Create(&region).Error; err == nil {
			regionMap[r.ID] = region.ID
		}
	}
	fmt.Println("✅ Regions Imported.")

	dataCity, err := os.ReadFile("assets/sql-moroccan-cities/json/ville.json")
	if err != nil {
		fmt.Println("❌ Error reading ville.json:", err)
		return
	}

	var jsonCities []JsonCity
	json.Unmarshal(dataCity, &jsonCities)

	var citiesToInsert []models.City
	for _, c := range jsonCities {
		realRegionID, exists := regionMap[c.RegionID]
		if exists {
			citiesToInsert = append(citiesToInsert, models.City{
				Ville:    c.Ville,
				RegionID: realRegionID,
			})
		}
	}

	if len(citiesToInsert) > 0 {
		db.CreateInBatches(citiesToInsert, 100)
	}
	fmt.Println("✅ Cities Imported Successfully.")
}

func SeedBanks(db *gorm.DB) {
	var count int64
	db.Model(&models.Bank{}).Count(&count)
	if count > 0 {
		fmt.Println("🏥 Banks already seeded.")
		return
	}

	fmt.Println("⏳ Seeding Hospitals Data...")

	data, err := os.ReadFile("assets/banks.json")
	if err != nil {
		fmt.Println("❌ Error reading banks.json:", err)
		return
	}

	var jsonBanks []JsonBank
	if err := json.Unmarshal(data, &jsonBanks); err != nil {
		fmt.Println("❌ Error parsing JSON:", err)
		return
	}

	for _, b := range jsonBanks {
		var stock []models.BloodGroup
		stockMap := make(map[string]int)
		for _, bg := range b.BloodGroups {
			stockMap[bg.Group] = bg.Value
		}

		allGroups := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
		for _, g := range allGroups {
			quantity := 0
			if val, ok := stockMap[g]; ok {
				quantity = val
			}
			stock = append(stock, models.BloodGroup{
				Group:    g,
				Quantity: quantity,
			})
		}
		    fmt.Printf("📥 Reading JSON: Name='%s', City='%s'\n", b.Name, )
		bank := models.Bank{
			Name:        b.Name,
			City: 		b.City,
			Email:       b.Email,
			PhoneNumber: b.PhoneNumber,
			Website:     b.Website,
			Address: models.Address{
				Area:       b.Address.Area,
				City:      b.Address.City,
				State:      b.Address.State,
				Country:    "Morocco",
				PostalCode: b.Address.PostalCode,
				Latitude:   b.Address.Latitude,
				Longitude:  b.Address.Longitude,
			},
			BloodGroups: stock,
		}

		if err := db.Create(&bank).Error; err != nil {
			fmt.Printf("⚠️ Failed to seed bank %s: %v\n", b.Name, err)
		}
	}
	fmt.Println("✅ Hospitals Seeded Successfully.")
}

func SeedRolesAndAdmin(db *gorm.DB) {
	roles := []string{"Admin", "Sponsor", "Donor", "Moderator"}
	for _, r := range roles {
		var count int64
		db.Model(&models.Role{}).Where("name = ?", r).Count(&count)
		if count == 0 {
			db.Create(&models.Role{Name: r})
			fmt.Printf("✅ Role Created: %s\n", r)
		}
	}

	var adminCount int64
	db.Model(&models.User{}).Where("user_name = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("younsse1234"), bcrypt.DefaultCost)

		admin := models.User{
			UserName:     "admin",
			Email:        "admin@lifeline.ma",
			PasswordHash: string(hash),
			Name:         "Super Admin",
			BloodGroup:   "O+",
		}

		db.Create(&admin)

		var adminRole models.Role
		db.Where("name = ?", "Admin").First(&adminRole)
		db.Model(&admin).Association("Roles").Append(&adminRole)

		fmt.Println(" Admin Account Created: admin / younsse1234")
	}
}

func SeedFakeUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("email LIKE ?", "%@fake.com").Count(&count)
	if count > 0 { return }

	fmt.Println("🌱 Seeding Fake Users...")
	
	pass, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	// Fetch Donor Role ONCE (Optimized)
	var donorRole models.Role
	db.Where("name = ?", "Donor").First(&donorRole)

	users := []models.User{
		{Name: "Karim O+", UserName: "karim", Email: "karim@fake.com",Photoprofile: "", BloodGroup: "O+", PasswordHash: string(pass)},
		{Name: "Fatima A-", UserName: "fatima", Email: "fatima@fake.com", Photoprofile: "",BloodGroup: "A-", PasswordHash: string(pass)},
		{Name: "Said AB+", UserName: "said", Email: "said@fake.com", Photoprofile: "",BloodGroup: "AB+", PasswordHash: string(pass)},
		{Name: "Anass AB+", UserName: "anass", Email: "anass@fake.com",Photoprofile: "", BloodGroup: "AB+", PasswordHash: string(pass)},
	}

	for _, u := range users {
		db.Create(&u)
		db.Model(&u).Association("Roles").Append(&donorRole)
	}
	fmt.Println("✅ Fake Users Created (Pass: 123456)")
}