package validation

import (
	user_model "openidea-idea-social-media-app/models/user"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterValidation(validator *validator.Validate) {
	validator.RegisterStructValidation(mustValidUserRequest, user_model.UserRegisterRequest{})
	validator.RegisterStructValidation(mustValidUserRequest, user_model.UserLoginRequest{})
}

func mustValidUserRequest(sl validator.StructLevel) {
	registerRequest := sl.Current().Interface().(user_model.UserRegisterRequest)

	if registerRequest.CredentialType == "email" {
		email := registerRequest.CredentialValue
		if !isValidEmail(email) {
			sl.ReportError(registerRequest.CredentialType, "CredentialValue", "CredentialValue", "credentialValue", "")
		}
	} else if registerRequest.CredentialType == "phone" {
		phone := registerRequest.CredentialValue
		if !isValidPhone(phone) {
			sl.ReportError(registerRequest.CredentialType, "CredentialValue", "CredentialValue", "credentialValue", "")
		}
	} else {
		sl.ReportError(registerRequest.CredentialType, "CredentialType", "CredentialType", "credentialType", "")
	}

}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(email)
}

func isValidPhone(phoneNumber string) bool {
	// Regular expression for phone number validation
	// ^\+(?:[0-9] ?){6,12}[0-9]$
	// Explanation:
	// ^               - Start of string
	// \+              - Literal '+'
	// (?:[0-9] ?)     - Non-capturing group for digits optionally followed by a space, repeated 6 to 12 times
	// [0-9]           - Last digit
	// $               - End of string
	re := regexp.MustCompile(`^\+(?:[0-9] ?){6,12}[0-9]$`)
	return re.MatchString(phoneNumber) && (len(phoneNumber) >= 7 && len(phoneNumber) <= 13)
}
