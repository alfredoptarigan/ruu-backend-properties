package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/helpers"
	"alfredo/ruu-properties/pkg/models"
)

type UserRepository interface {
	Login(email string, password string) (user models.User, err error)
	Register(request dtos.UserRegisterRequest) error
	FindUserByEmail(email string) (models.User, error)
	FindUserByUuid(uuid string) (models.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// FindUserByUuid implements UserRepository.
func (u *userRepositoryImpl) FindUserByUuid(uuid string) (models.User, error) {
	var user models.User

	if err := u.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error retrieving user: %w", err)
	}
	return user, nil
}

func (u *userRepositoryImpl) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := u.db.Unscoped().Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error retrieving user: %w", err)
	}
	return user, nil
}

func (u *userRepositoryImpl) Register(request dtos.UserRegisterRequest) error {
	var existingUser models.User
	if err := u.db.Unscoped().Debug().Where("email = ?", request.Email).First(&existingUser).Error; err != nil && existingUser.UUID != "" {
		return fmt.Errorf("%s", "Email already exists on database")
	}

	if err := u.db.Unscoped().Where("phone_number = ?", request.PhoneNumber).First(&existingUser).Error; err != nil && existingUser.UUID != "" {
		return fmt.Errorf("%s", "Phone number already exists on database")
	}

	// Hash the password using argon2
	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Name:        request.Name,
			Email:       request.Email,
			Password:    hashedPassword, // Store the hashed password
			PhoneNumber: request.PhoneNumber,
			Role:        request.Role,
			Image:       request.Image,
		}

		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) Login(email string, password string) (user models.User, err error) {
	if err := u.db.Where("email = ? AND password = ? AND deleted_at IS NULL", email, password).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		} else {
			// Handle other errors
			panic(err)
		}
	}

	return user, nil
}
