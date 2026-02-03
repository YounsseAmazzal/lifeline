package services

import (
	"lifeline/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService interface (باش نحافظو على clean architecture)
type TokenService interface {
	CreateToken(user *models.User, existingToken string) (string, error)
}

type tokenService struct {
	secretKey []byte
}

// NewTokenService constructor
func NewTokenService() TokenService {
	// كنجيبو المفتاح السري من .env
	secret := os.Getenv("TOKEN_KEY")
	if secret == "" {
		secret = "super_secret_default_key_change_me" // Fallback (غير للتجربة)
	}
	return &tokenService{
		secretKey: []byte(secret),
	}
}

// CreateToken implementation
func (s *tokenService) CreateToken(user *models.User, existingToken string) (string, error) {
	// 1. Define Claims
	claims := jwt.MapClaims{
		"nameid":      user.ID,       // UserID
		"unique_name": user.UserName, // UserName
	}

	// 2. Add Roles (إيلا كانو ديجا معمرين فاليوزر)
	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	if len(roles) > 0 {
		claims["role"] = roles
	}

	// 3. Expiration Logic
	// (بسطت اللوجيك ديال existingToken حيت فـ Go كنفضلوا نجددو التوكن ديريكت)
	expireTime := time.Now().Add(24 * time.Hour) // صالح لـ 24 ساعة
	
	// إيلا بغيتي تطبق اللوجيك ديال existingToken بالضبط:
	if existingToken != "" {
		// Parsing existing token to get exp requires extra logic, 
		// usually for refresh tokens. For MVP, fresh token is safer.
	}
	
	claims["exp"] = expireTime.Unix()

	// 4. Create & Sign Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // استعملنا HS512 بحال C#
	
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}