package requests

type RegisterAccountRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required"`
	ConfirmationPassword string `json:"confirmationPassword" validate:"required,eqfield=Password"`
	IsCompany            int    `json:"isCompany" validate:"oneof=0 1"`
}

type LoginAccountRequest struct {
	Email    string `json:"email" validate:"required,email,min=99"`
	Password string `json:"password" validate:"required"`
}
