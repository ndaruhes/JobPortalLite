package repositories

import (
	"errors"
	"job-portal-lite/domain/user/entities"
	"job-portal-lite/shared/databases"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() *userRepository {
	return &userRepository{
		db: databases.Connect(),
	}
}

func (u userRepository) RegisterUser(user *entities.User) (*entities.User, error) {
	// Check if user already exists
	var count int64
	err := u.db.Where("email = ?", user.Email).Find(&entities.User{}).Count(&count).Error
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("You are already registered")
	}

	// Create User if already exists
	err = u.db.Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u userRepository) LoginUser(user *entities.User) (*entities.User, error) {
	var userResponse *entities.User
	err := u.db.Where("email = ?", user.Email).First(&userResponse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("You are not registered")
		}
		return nil, err
	}

	return userResponse, err
}

func (u userRepository) GetCurrentUser(user *entities.User) (*entities.User, error) {
	var userResponse *entities.User
	err := u.db.Where("email = ?", user.Email).First(&userResponse).Error
	return userResponse, err
}

func (u userRepository) GetUserById(id int) (*entities.User, error) {
	var userResponse *entities.User
	err := u.db.Where("id = ?", id).First(&userResponse).Error
	return userResponse, err
}
