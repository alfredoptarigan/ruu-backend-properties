package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"alfredo/ruu-properties/pkg/injectors"
	"alfredo/ruu-properties/pkg/router"
)

// Import swagger

// @title RUU Properties API
// @version 1.0
// @description API documentation for RUU Properties Backend
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email your-email@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9090
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	time.Local = time.UTC

	server := injectors.InitializeApplication()
	server.RegisterMiddlewares()
	router.InitializeRouterV1(server)

	go func() {
		if err := server.App.Listen(fmt.Sprintf(":%s", server.Config.GetString("application.port"))); err != nil {
			server.Logger.Error("Failed to start server", "error", nil)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.Logger.Info("Shutting down server")
	if err := server.App.ShutdownWithTimeout(10 * time.Second); err != nil {
		server.Logger.Error("Failed to shutdown server", "error", err)
	}

	server.Logger.Info("Server shut down successfully")
	defer server.CloseConnectionDatabase()
}
