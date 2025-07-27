package admin

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/models"
)

func IsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil data user dari context yang sudah di-set oleh JwtMiddleware
		userData := c.Locals("user")
		if userData == nil {
			log.Println("Error: User data not found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
				Message: "Unauthorized",
				Code:    fiber.StatusUnauthorized,
			})
		}

		// Type assertion untuk mendapatkan user object
		user, ok := userData.(*models.User)
		if !ok {
			log.Println("Error: Invalid user data type")
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
				Message: "Internal Server Error",
				Code:    fiber.StatusInternalServerError,
			})
		}

		// Periksa apakah user memiliki role admin
		if user.Role != "admin" {
			log.Printf("Access denied for user %s with role %s", user.Email, user.Role)
			return c.Status(fiber.StatusForbidden).JSON(dtos.ErrorResponseDTO{
				Message: "Access denied. Admin role required",
				Code:    fiber.StatusForbidden,
			})
		}

		// Jika user adalah admin, lanjutkan ke handler berikutnya
		return c.Next()
	}
}
