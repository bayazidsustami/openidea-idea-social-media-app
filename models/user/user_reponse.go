package user_model

type UserRegisterResponse[T UserData] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type UserData interface {
	GetName() string
}

type UserEmailDataResponse struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (ue *UserEmailDataResponse) GetName() string {
	return ue.Name
}

type UserPhoneDataResponse struct {
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (up *UserPhoneDataResponse) GetName() string {
	return up.Name
}

type UserLoginResponse struct {
	Message string                     `json:"message"`
	Data    UserEmailPhoneDataResponse `json:"data"`
}

type UserEmailPhoneDataResponse struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
