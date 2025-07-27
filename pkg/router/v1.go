package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"alfredo/ruu-properties/config"
	"alfredo/ruu-properties/pkg/injectors"

	_ "alfredo/ruu-properties/docs" // Import ini penting untuk swagger
)

func InitializeRouterV1(server *config.Application) {
	// Gunakan handler untuk path /swagger/* bukan middleware
	server.App.Get("/swagger/*", swagger.HandlerDefault) // Opsi sederhana

	// ATAU gunakan konfigurasi yang lebih spesifik
	server.App.Get("/swagger/*", swagger.New(swagger.Config{
		Title:        "RUU Properties API Documentation",
		URL:          "/swagger/doc.json", // Pastikan ini sesuai dengan path yang benar
		DeepLinking:  true,
		Layout:       "BaseLayout",
		DocExpansion: "none",
	}))
	api := server.App.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Ping godoc
			// @Summary Health check endpoint
			// @Description Get server status
			// @Tags system
			// @Accept json
			// @Produce json
			// @Success 200 {object} map[string]string
			// @Router /ping [get]
			v1.Get("/ping", func(ctx *fiber.Ctx) error {
				return ctx.JSON(fiber.Map{
					"message": "pong",
				})
			})

			auth := v1.Group("/auth")
			{
				authController := injectors.InitializeAuthController()
				authController.Router(auth)
			}

			user := v1.Group("/user")
			{
				userController := injectors.InitializeUserController()
				userController.Router(user)

			}

		}

	}
}

// SetupRoutes sets up all routes for testing purposes
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Ping endpoint
			v1.Get("/ping", func(ctx *fiber.Ctx) error {
				return ctx.JSON(fiber.Map{
					"message": "pong",
				})
			})

			auth := v1.Group("/auth")
			{
				authController := injectors.InitializeAuthController()
				authController.Router(auth)
			}

			user := v1.Group("/user")
			{
				userController := injectors.InitializeUserController()
				userController.Router(user)
			}
		}
	}
}
