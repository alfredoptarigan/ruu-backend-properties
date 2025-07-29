//go:build wireinject
// +build wireinject

package injectors

import (
	"github.com/google/wire"

	"alfredo/ruu-properties/config"
	"alfredo/ruu-properties/pkg/controllers"
	"alfredo/ruu-properties/pkg/repositories"
	"alfredo/ruu-properties/pkg/services"
	"alfredo/ruu-properties/pkg/validator"
)

var initDBPostgresSet = wire.NewSet(
	config.InitDatabasePostgres,
)

var redisSet = wire.NewSet(
	config.InitRedis,
	repositories.NewRedisRepository,
	services.NewRedisService,
)

var jwtSet = wire.NewSet(
	services.NewJwtService,
)

var authSet = wire.NewSet(
	redisSet,
	initDBPostgresSet,
	services.NewUserService,
	repositories.NewUserRepository,
	validator.NewValidator,
)

func InitializeApplication() *config.Application {
	wire.Build(config.NewApplication, config.InitDatabasePostgres)
	return nil
}

func InitializeAuthController() controllers.AuthController {
	wire.Build(
		authSet,
		jwtSet,
		controllers.NewAuthController,
		services.NewAuthService,
	)

	return nil
}

func InitializeUserController() controllers.UserController {
	wire.Build(
		authSet,
		controllers.NewUserController,
	)

	return nil
}

func InitializeClientController() controllers.ClientController {
	wire.Build(
		authSet,
		controllers.NewClientController,
		services.NewClientService,
		repositories.NewClientRepository,
	)

	return nil
}

func InitializeFeatureController() controllers.FeatureController {
	wire.Build(
		authSet,
		controllers.NewFeatureController,
		services.NewFeatureService,
		repositories.NewFeatureRepository,
	)

	return nil
}
