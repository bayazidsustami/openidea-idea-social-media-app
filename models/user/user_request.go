package user_model

type UserRegisterRequest struct {
	CredentialType  string `json:"credentialType" validate:"required,oneof=email phone"`
	CredentialValue string `json:"credentialValue" validate:"required"`
	Name            string `json:"name" validate:"required,min=5,max=50"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type UserLoginRequest struct {
	CredentialType  string `json:"credentialType" validate:"required,oneof=email phone"`
	CredentialValue string `json:"credentialValue" validate:"required"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type UpdatePhoneRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type UpdateAccountRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	ImageUrl string `json:"imageUrl" validate:"required,imageurl"`
}
