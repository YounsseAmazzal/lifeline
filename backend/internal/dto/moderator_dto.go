package dto
type MemberResponse struct {
	ID          uint   `json:"id"`
	UserName    string `json:"username"`
	Name        string `json:"name"`
	PhotoURL    string `json:"photo_url"`
	Gender      string `json:"gender"`
	City        string `json:"city"`
	BloodGroup  string `json:"blood_group"`
	LastActive  string `json:"last_active"`
	Available   bool   `json:"available"`
}