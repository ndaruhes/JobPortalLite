package responses

type LoginResponse struct {
	Data    TokenResponse `json:"data"`
	Message string        `json:"message"`
	Status  string        `json:"status"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type CurrentUserResponse struct {
	Data    UserResponse `json:"data"`
	Message string       `json:"message"`
	Status  string       `json:"status"`
}

type RegisterResponse struct {
	Status string `json:"status"`
}

type TokenDecoded struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
