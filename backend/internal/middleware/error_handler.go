package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"` 
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("❌ Error: %v\n", err)

	return c.Status(code).JSON(APIError{
		StatusCode: code,
		Message:    message,
	})
}