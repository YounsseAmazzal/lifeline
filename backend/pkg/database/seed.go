package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lifeline/internal/models"

	"gorm.io/gorm"
)

type JsonRegion struct {
	ID     string `json:"id"` 
	Region string `json:"region"`
}

type JsonCity struct {
	ID       string `json:"id"`
	Ville    string `json:"ville"`
	RegionID string `json:"region"` 
}

func SeedMoroccanCities(db *gorm.DB) {
	var count int64
	db.Model(&models.Region{}).Count(&count)
	if count > 0 {
		fmt.Println("🇲🇦 Moroccan Regions & Cities already seeded.")
		return
	}

	fmt.Println("⏳ Seeding Moroccan Data...")

	dataRegion, _ := ioutil.ReadFile("assets/sql-moroccan-cities/json/region.json")
	var jsonRegions []JsonRegion
	json.Unmarshal(dataRegion, &jsonRegions)

	for _, r := range jsonRegions {
		db.Create(&models.Region{
			Region: r.Region,
		})
	}
	fmt.Println(" Regions Imported.")

	dataCity, _ := ioutil.ReadFile("assets/sql-moroccan-cities/json/ville.json")
	var jsonCities []JsonCity
	json.Unmarshal(dataCity, &jsonCities)


	for _, c := range jsonCities {
		var region models.Region
		db.Where("id = ?", c.RegionID).First(&region)

		if region.ID != 0 {
			db.Create(&models.City{
				Ville:    c.Ville,
				RegionID: region.ID,
			})
		}
	}
	fmt.Println("✅ Cities Imported.")
}

type JsonAddress struct {
	Area       string `json:"Area"`
	City       string `json:"City"`
	State      string `json:"State"`
	Country    string `json:"Country"`
	PostalCode string `json:"PostalCode"`
}

type JsonBloodGroup struct {
	Group string `json:"Group"`
	Value int    `json:"Value"` 
}

type JsonBank struct {
	Name        string           `json:"Name"`
	PhoneNumber string           `json:"PhoneNumber"`
	Email       string           `json:"Email"`
	Website     string           `json:"Website"`
	Address     JsonAddress      `json:"Address"`
	BloodGroups []JsonBloodGroup `json:"BloodGroups"`
}

func SeedBanks(db *gorm.DB) {
	var count int64
	db.Model(&models.Bank{}).Count(&count)
	if count > 0 {
		fmt.Println(" Banks already seeded.")
		return
	}

	fmt.Println("⏳ Seeding Hospitals Data...")

	data, err := ioutil.ReadFile("assets/banks.json")
	if err != nil {
		fmt.Println(" Error reading banks.json:", err)
		return
	}

	var jsonBanks []JsonBank
	err = json.Unmarshal(data, &jsonBanks)
	if err != nil {
		fmt.Println(" Error parsing JSON:", err)
		return
	}

	for _, b := range jsonBanks {
		addr := models.Address{
			Area:       b.Address.Area,
			City:       b.Address.City,
			State:      b.Address.State,
			Country:    "Morocco", 
			PostalCode: b.Address.PostalCode,
		}

		var stock []models.BloodGroup
		if len(b.BloodGroups) > 0 {
			for _, bg := range b.BloodGroups {
				stock = append(stock, models.BloodGroup{
					Group:    bg.Group,
					Quantity: bg.Value, 
				})
			}
		} else {
			groups := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
			for _, g := range groups {
				stock = append(stock, models.BloodGroup{Group: g, Quantity: 0})
			}
		}

		// Create Bank
		bank := models.Bank{
			Name:        b.Name,
			Email:       b.Email,
			PhoneNumber: b.PhoneNumber,
			Website:     b.Website,
			Address:     addr,
			BloodGroups: stock,
		}

		if err := db.Create(&bank).Error; err != nil {
			fmt.Printf(" Failed to seed bank %s: %v\n", b.Name, err)
		}
	}
	fmt.Println(" Hospitals Seeded Successfully.")
}
