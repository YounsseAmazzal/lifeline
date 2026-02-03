package dto

// UserResponse: 
type UserResponse struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	PhotoURL string `json:"photo_url"` 
}