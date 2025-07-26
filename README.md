# RUU Properties Backend

This is the backend service for the RUU Properties application, built with Go, Fiber, GORM, and PostgreSQL. It provides APIs for user management, authentication, and property-related functionalities.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Database Operations](#database-operations)
- [Testing](#testing)
- [Linting](#linting)
- [Dependency Injection](#dependency-injection)
- [Project Structure](#project-structure)
- [API Documentation (Swagger)](#api-documentation-swagger)

## Features

- User Registration and Authentication (JWT)
- Image Uploads
- Database Migrations and Seeding
- Dependency Injection with Wire
- Structured Logging
- Redis Caching (for future use/session management)

## Technologies Used

- **Go**: Programming Language
- **Fiber**: Fast, Express-inspired web framework
- **GORM**: ORM library for Go
- **PostgreSQL**: Relational Database
- **Redis**: In-memory data structure store
- **go-playground/validator**: Request validation
- **golang-jwt/jwt**: JWT implementation
- **google/uuid**: UUID generation
- **google/wire**: Dependency Injection
- **spf13/viper**: Configuration management
- **golang.org/x/crypto**: Cryptographic functions (e.g., Argon2 for password hashing)

## Getting Started

Follow these instructions to set up and run the project locally.

### Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.24.5 or higher)
- PostgreSQL
- Redis
- Git
- `golangci-lint` (for linting): `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- `wire` (for dependency injection): `go install github.com/google/wire/cmd/wire@latest`
- `swag` (for Swagger documentation): `go install github.com/swaggo/swag/cmd/swag@latest`

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd ruu-properties-backend
   ```

2. **Download Go modules:**
   ```bash
   go mod tidy
   ```

### Configuration

Create a `.env` file in the project root based on `example.config.yaml` and fill in your database and other environment variables.

```yaml:example.config.yaml
# Example .env content
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your_jwt_secret
JWT_EXPIRATION_MINUTES=60

APP_PORT=8080
```

## Running the Application

To build and run the application:

1. **Generate Wire files (Dependency Injection):**
   ```bash
   make wire
   ```

2. **Build the application:**
   ```bash
   make build
   ```
   This will create an executable `gic-crm` in the `bin/` directory.

3. **Run the application:**
   ```bash
   make run
   ```
   Alternatively, you can run the executable directly:
   ```bash
   ./bin/gic-crm
   ```

## Database Operations

- **Run migrations:**
  ```bash
  make migrate
  ```

- **Seed the database:**
  ```bash
  make seed
  ```

## Testing

To run all tests:

```bash
make test
```

## Linting

To run the linter:

```bash
make lint
```

## Dependency Injection

This project uses `google/wire` for dependency injection. After adding new dependencies or modifying the dependency graph, you need to regenerate the `wire_gen.go` files:

```bash
make wire
```

## Project Structure