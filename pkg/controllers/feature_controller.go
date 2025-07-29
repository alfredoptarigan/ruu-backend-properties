package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/helpers"
	"alfredo/ruu-properties/pkg/middleware/jwt"
	"alfredo/ruu-properties/pkg/services"
)

type FeatureController interface {
	Create(c *fiber.Ctx) error
	Router(router fiber.Router)
}

type featureControllerImpl struct {
	featureService services.FeatureService
	userService    services.UserService
	redisService   services.RedisService
}

// Router implements FeatureController.
func (f *featureControllerImpl) Router(router fiber.Router) {
	withMiddleware := router.Use(jwt.JwtMiddleware(f.userService, f.redisService))
	{
		withMiddleware.Post("/", f.Create)
	}
}

// Create implements FeatureController.
func (f *featureControllerImpl) Create(c *fiber.Ctx) error {
	var request dtos.FeatureRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
			Errors:  []string{err.Error()},
		})
	}

	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
			Errors:  helpers.FormatValidationError(err),
		})
	}

	feature, err := f.featureService.Create(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
			Message: "Failed to create feature",
			Code:    fiber.StatusInternalServerError,
			Errors:  []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dtos.SuccessResponse{
		Message: "Feature created successfully",
		Data:    feature,
		Success: true,
	})
}

func NewFeatureController(featureService services.FeatureService, userService services.UserService, redisService services.RedisService) FeatureController {
	return &featureControllerImpl{featureService: featureService, userService: userService, redisService: redisService}
}
