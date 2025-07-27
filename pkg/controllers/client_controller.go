package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/helpers"
	"alfredo/ruu-properties/pkg/middleware/jwt"
	"alfredo/ruu-properties/pkg/services"
)

type ClientController interface {
	Create(c *fiber.Ctx) error
	Router(router fiber.Router)
}

type clientControllerImpl struct {
	redisService  services.RedisService
	userService   services.UserService
	clientService services.ClientService
}

// Router implements ClientController.
func (c *clientControllerImpl) Router(router fiber.Router) {
	withMiddleware := router.Use(jwt.JwtMiddleware(c.userService, c.redisService))
	{
		withMiddleware.Post("/", c.Create)
	}
}

// Create Client godoc
// @Summary Create a new client
// @Description Create a new client for renting/selling properties
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dtos.ClientRequest true "Client request"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.ClientResponse}
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 401 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /clients [post]
func (cs *clientControllerImpl) Create(c *fiber.Ctx) error {
	var request dtos.ClientRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Message: "Invalid request body",
		})
	}

	// Validate the request first
	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
			Errors:  helpers.FormatValidationError(err),
		})
	}

	response, err := cs.clientService.Create(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	return c.JSON(dtos.SuccessResponse{
		Data: response,
	})
}

func NewClientController(
	redisService services.RedisService,
	userService services.UserService,
	clientService services.ClientService,
) ClientController {
	return &clientControllerImpl{
		redisService:  redisService,
		userService:   userService,
		clientService: clientService,
	}
}
