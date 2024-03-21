package validation

import (
	"net/url"
	user_model "openidea-idea-social-media-app/models/user"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterValidation(validator *validator.Validate) {
	validator.RegisterStructValidation(mustValidRegisterRequest, user_model.UserRegisterRequest{})
	validator.RegisterStructValidation(mustValidLoginRequest, user_model.UserLoginRequest{})
	validator.RegisterStructValidation(mustValidEmailRequest, user_model.UpdateEmailRequest{})
	validator.RegisterStructValidation(mustValidPhoneRequest, user_model.UpdatePhoneRequest{})
	validator.RegisterValidation("imageurl", mustValidImageUrl)
}

func mustValidRegisterRequest(sl validator.StructLevel) {
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

func mustValidLoginRequest(sl validator.StructLevel) {
	registerRequest := sl.Current().Interface().(user_model.UserLoginRequest)

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

func mustValidEmailRequest(sl validator.StructLevel) {
	emailReq := sl.Current().Interface().(user_model.UpdateEmailRequest)

	if !isValidEmail(emailReq.Email) {
		sl.ReportError(emailReq.Email, "Email", "Email", "email", "")
	}
}

func mustValidPhoneRequest(sl validator.StructLevel) {
	phoneReq := sl.Current().Interface().(user_model.UpdatePhoneRequest)

	if !isValidPhone(phoneReq.Phone) {
		sl.ReportError(phoneReq.Phone, "Phone", "Phone", "phone", "")
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

func mustValidImageUrl(fl validator.FieldLevel) bool {
	urlString := fl.Field().String()

	// Parse the URL
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}

	// Get the file extension
	parts := strings.Split(u.Path, ".")
	extension := parts[len(parts)-1]

	// Check if the extension is jpg or jpeg
	if extension != "jpg" && extension != "jpeg" {
		return false
	}

	return true
}
