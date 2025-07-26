package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/services"
)

type UserController interface {
	Router(router fiber.Router)
	Register(c *fiber.Ctx) error
}

type userControllerImpl struct {
	redisService services.RedisService
	userService  services.UserService
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "User's full name"
// @Param email formData string true "User's email address"
// @Param password formData string true "Password (min 6 characters)"
// @Param confirmation_password formData string true "Confirm password"
// @Param phone_number formData string true "User's phone number"
// @Param role formData string true "User role"
// @Param image formData file false "User profile image"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.UserRegisterRequest}
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /user/register [post]
func (u *userControllerImpl) Register(c *fiber.Ctx) error {
	var request dtos.UserRegisterRequest

	// Parse form fields
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
			Errors:  err.Error(),
		})
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll("./uploads", 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Failed to create upload directory",
				Code:    fiber.StatusInternalServerError,
				Errors:  err.Error(),
			})
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
		filepath := fmt.Sprintf("./uploads/%s", filename)

		// Save the file
		if err := c.SaveFile(file, filepath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Failed to save file",
				Code:    fiber.StatusInternalServerError,
				Errors:  err.Error(),
			})
		}

		// Set the image path in the request
		request.Image = filepath
	}

	if err := u.userService.Register(request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    nil,
	})
}

// Router implements UserController.
func (u *userControllerImpl) Router(router fiber.Router) {
	router.Post("/register", u.Register)
}

func NewUserController(
	redisService services.RedisService,
	userService services.UserService,
) UserController {
	return &userControllerImpl{
		redisService: redisService,
		userService:  userService,
	}
}
