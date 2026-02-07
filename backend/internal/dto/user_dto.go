package dto

// UserResponse (Zidna Role)
type UserResponse struct {
	UserName string `json:"userName"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	PhotoURL string `json:"photoUrl"`
	Token    string `json:"token"`
	Role     string `json:"role"` // <--- ZID HADI
}