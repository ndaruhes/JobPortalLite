package interfaces

import (
	"job-portal-lite/domain/user/entities"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"
)

type AuthUseCase interface {
	RegisterUser(request *requests.RegisterAccountRequest) (*responses.RegisterResponse, error)
	LoginUser(request *requests.LoginAccountRequest) (*responses.LoginResponse, error)
	GetCurrentUser(*entities.User) (*responses.CurrentUserResponse, error)
}

type UserRepository interface {
	RegisterUser(*entities.User) (*entities.User, error)
	LoginUser(user *entities.User) (*entities.User, error)
	GetCurrentUser(*entities.User) (*entities.User, error)
	GetUserById(int) (*entities.User, error)
}
