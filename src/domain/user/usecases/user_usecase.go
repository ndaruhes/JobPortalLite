package usecases

import (
	"errors"
	"job-portal-lite/domain/user/entities"
	"job-portal-lite/domain/user/interfaces"
	"job-portal-lite/domain/user/repositories"

	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"
	"job-portal-lite/shared/helpers"
)

type userUseCase struct {
	repo interfaces.UserRepository
}

func NewUserUseCase() *userUseCase {
	return &userUseCase{
		repo: repositories.NewUserRepository(),
	}
}

func (u userUseCase) RegisterUser(request *requests.RegisterAccountRequest) (*responses.RegisterResponse, error) {
	_, err := u.repo.RegisterUser(&entities.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: helpers.HashPassword(request.Password),
		Role: func() string {
			if request.IsCompany == 1 {
				return "Company"
			}
			return "Candidate"
		}(),
	})

	if err != nil {
		return nil, err
	}

	response := &responses.RegisterResponse{
		Status: "OK",
	}

	return response, err
}

func (u userUseCase) LoginUser(request *requests.LoginAccountRequest) (*responses.LoginResponse, error) {
	userLogin, err := u.repo.LoginUser(&entities.User{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		response := &responses.LoginResponse{
			Message: err.Error(),
		}

		return response, err
	}

	checkPassword := helpers.VerifyPassword([]byte(userLogin.Password), []byte(request.Password))
	if !checkPassword {
		return nil, errors.New("Wrong Credentials")
	}

	token, err := helpers.GenerateToken(userLogin.ID, userLogin.Email, userLogin.Name, userLogin.Role)
	if err != nil {
		response := &responses.LoginResponse{
			Message: err.Error(),
		}

		return response, err
	}

	response := &responses.LoginResponse{
		Data: responses.TokenResponse{
			Token: token,
		},
		Message: "Login Succesfully",
		Status:  "OK",
	}

	return response, err
}

func (u userUseCase) GetCurrentUser(user *entities.User) (*responses.CurrentUserResponse, error) {
	data, err := u.repo.GetCurrentUser(user)

	if err != nil {
		return nil, err
	}

	return &responses.CurrentUserResponse{
		Data: responses.UserResponse{
			Name: data.Name,
			Role: data.Role,
		},

		Message: "Profile Showed",
		Status:  "OK",
	}, nil
}
