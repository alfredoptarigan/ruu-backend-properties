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
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Router(router fiber.Router)
}

type clientControllerImpl struct {
	redisService  services.RedisService
	userService   services.UserService
	clientService services.ClientService
}

// Update Client godoc
// @Summary Update an existing client
// @Description Update client information by ID
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Client ID"
// @Param request body dtos.ClientRequest true "Client request"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.ClientResponse}
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 401 {object} dtos.ErrorResponseDTO
// @Failure 404 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /clients/{id}/update [put]
func (cs *clientControllerImpl) Update(c *fiber.Ctx) error {
	var request dtos.ClientUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid request body",
			Errors:  []string{err.Error()},
		})
	}

	uuid := c.Params("id")
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid client ID",
		})
	}
	if !helpers.CheckLengthUUID(uuid) {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid client ID",
		})
	}
	request.UUID = uuid

	clientResponse, err := cs.clientService.Update(request)
	if err != nil {
		if err.Error() == "client not found" {
			return c.Status(fiber.StatusNotFound).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: err.Error(),
				Errors:  []string{err.Error()},
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: err.Error(),
				Errors:  []string{err.Error()},
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(dtos.SuccessResponse{
		Success: true,
		Message: "Client updated successfully",
		Data:    clientResponse,
	})
}

// Delete Client godoc
// @Summary Delete a client
// @Description Delete a client by ID
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Client ID"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponseDTO
// @Failure 404 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /clients/{id}/delete [delete]
func (cs *clientControllerImpl) Delete(c *fiber.Ctx) error {
	uuid := c.Params("id")
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid client ID",
		})
	}

	if !helpers.CheckLengthUUID(uuid) {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid client ID",
		})
	}

	if err := cs.clientService.Delete(uuid); err != nil {
		if err.Error() == "client not found" {
			return c.Status(fiber.StatusNotFound).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: err.Error(),
				Errors:  []string{err.Error()},
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: err.Error(),
				Errors:  []string{err.Error()},
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(dtos.SuccessResponse{
		Success: true,
		Message: "Client deleted successfully",
	})
}

// GetAll Client godoc
// @Summary Get all clients
// @Description Get a list of all clients with pagination and search functionality
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param search query string false "Search term"
// @Param search_by query string false "Field to search by (name, email, phone_number, contact_person)" default(name)
// @Param sort_by query string false "Field to sort by (name, email, created_at, updated_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} dtos.PaginatedSuccessResponse{data=[]dtos.ClientResponse,meta=dtos.PaginationMeta}
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 401 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /clients [get]
func (cs *clientControllerImpl) GetAll(c *fiber.Ctx) error {
	var request dtos.ClientGetRequest
	if err := c.QueryParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid query parameters",
			Errors:  err.Error(),
		})
	}

	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}

	// Validasi tambahan untuk parameter search_by dan sort_by
	if request.SearchBy != "" {
		allowedSearchFields := map[string]bool{
			"name":           true,
			"email":          true,
			"phone_number":   true,
			"contact_person": true,
		}
		if !allowedSearchFields[request.SearchBy] {
			return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Invalid search_by parameter. Allowed values: name, email, phone_number, contact_person",
			})
		}
	}

	if request.SortBy != "" {
		allowedSortFields := map[string]bool{
			"name":       true,
			"email":      true,
			"created_at": true,
			"updated_at": true,
		}
		if !allowedSortFields[request.SortBy] {
			return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Invalid sort_by parameter. Allowed values: name, email, created_at, updated_at",
			})
		}
	}

	if request.SortOrder != "" && request.SortOrder != "asc" && request.SortOrder != "desc" {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid sort_order parameter. Allowed values: asc, desc",
		})
	}

	clients, paginationMeta, err := cs.clientService.GetAll(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Failed to fetch clients",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.PaginatedSuccessResponse{
		Success: true,
		Message: "Successfully fetched clients",
		Data:    clients,
		Meta:    *paginationMeta,
	})
}

// GetByID Client godoc
// @Summary Get a client by ID
// @Description Get detailed information of a specific client
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Client ID"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.ClientResponse}
// @Failure 401 {object} dtos.ErrorResponseDTO
// @Failure 404 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /clients/{id} [get]
func (cs *clientControllerImpl) GetByID(c *fiber.Ctx) error {
	uuid := c.Params("id")
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Client UUID is required",
			})
	}

	if !helpers.CheckLengthUUID(uuid) {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid UUID format.",
		})
	}

	client, err := cs.clientService.GetByID(uuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: err.Error(),
				Errors:  err,
			})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.SuccessResponse{
		Success: true,
		Message: "Successfully fetched client",
		Data:    client,
	})
}

// Router implements ClientController.
func (c *clientControllerImpl) Router(router fiber.Router) {
	withMiddleware := router.Use(jwt.JwtMiddleware(c.userService, c.redisService))
	{
		withMiddleware.Get("/", c.GetAll)
		withMiddleware.Get("/:id", c.GetByID)
		withMiddleware.Post("/", c.Create)
		withMiddleware.Put("/:id/update", c.Update)
		withMiddleware.Delete("/:id/delete", c.Delete)
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
		Success: true,
		Message: "Successfully created client",
		Data:    response,
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
