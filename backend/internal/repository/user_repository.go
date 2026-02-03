package repository

import (
	"lifeline/internal/dto"
	"lifeline/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Address").Preload("Photo").
		Where("user_name = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUsers(params dto.UserParams) ([]models.User, dto.PaginationHeader, error) {
	var users []models.User
	var totalCount int64

	query := r.db.Model(&models.User{}).Preload("Address").Preload("Photo")

	if params.CurrentUserName != "" {
		query = query.Where("user_name != ?", params.CurrentUserName)
	}

	if params.Gender != "" {
		query = query.Where("gender = ?", params.Gender)
	}

	if params.BloodGroup != "" {
		query = query.Where("blood_group = ?", params.BloodGroup)
	}

	now := time.Now()
	
	if params.MinAge > 0 {
		minDate := now.AddDate(-params.MinAge, 0, 0)
		query = query.Where("date_of_birth <= ?", minDate)
	}
	
	if params.MaxAge > 0 {
		maxDate := now.AddDate(-params.MaxAge-1, 0, 0) 
		query = query.Where("date_of_birth > ?", maxDate)
	}

	switch params.OrderBy {
	case "created":
		query = query.Order("created_at desc")
	default:
		query = query.Order("last_active desc")  
	}

	query.Count(&totalCount) 
	offset := (params.PageNumber - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&users).Error

	totalPages := int(totalCount) / params.PageSize
	if int(totalCount)%params.PageSize != 0 {
		totalPages++
	}

	paginationHeader := dto.PaginationHeader{
		PageNumber: params.PageNumber,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}

	return users, paginationHeader, err
}

func (r *UserRepository) LogUserActive(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("last_active", time.Now()).Error
}

func (r *UserRepository) UpdateProfile(user *models.User) error {
	return r.db.Save(user).Error
}