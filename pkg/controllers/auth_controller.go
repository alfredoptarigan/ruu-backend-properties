package controllers

import (
	"github.com/gofiber/fiber/v2"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/services"
)

type AuthController interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Router(router fiber.Router)
}

type authControllerImpl struct {
	authService  services.AuthService
	redisService services.RedisService
	jwtService   services.JwtService
}

func (a *authControllerImpl) Router(router fiber.Router) {
	router.Post("/login", a.Login)
	router.Post("/logout", a.Logout)
	router.Post("/refresh-token", a.RefreshToken)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.LoginResponse}
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /auth/login [post]
func (a *authControllerImpl) Login(c *fiber.Ctx) error {
	var request dtos.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
			Errors:  err.Error(),
		})
	}

	response, err := a.authService.Login(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dtos.ErrorResponseDTO{
				Success: false,
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
				Errors:  err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// / Logout godoc
// @Summary User logout
// @Description Invalidate user's refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /auth/logout [post]
func (a *authControllerImpl) Logout(c *fiber.Ctx) error {
	panic("unimplemented")
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dtos.SuccessResponse{data=dtos.GenerateTokenResponse}
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /auth/refresh-token [post]
func (a *authControllerImpl) RefreshToken(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewAuthController(
	authService services.AuthService,
	redisService services.RedisService,
	jwtService services.JwtService,
) AuthController {
	return &authControllerImpl{
		authService:  authService,
		redisService: redisService,
		jwtService:   jwtService,
	}
}
