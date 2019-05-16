package api

type MessageID string

// Get message returns message by it's id
func GetMessage(id MessageID) string {
	return messages[id]
}

var messages = map[MessageID]string{
	ErrorAPIInternalServerError:              "Internal server error",
	ErrorAPIUserUnableToCreate:               "Unable to create use. Try again or contact support team",
	ErrorAPIUserNotFound:                     "Specified user not found",
	ErrorAPIUsernameTaken:                    "Username is taken. Try to choose another one",
	ErrorAPIEmailTaken:                       "Email is taken. Try to choose another one",
	ErrorAPIInviteTokenServerError:           "Unable to create invite token. Try again or contact support team",
	ErrorAPIEmailNotSent:                     "Unable to send email. Try again or contact support team",
	ErrorAPIRequestPasswordWeak:              "Password is not strong enough",
	ErrorAPIRequestIncorrectEmailOrPassword:  "Incorrect email or password",
	ErrorAPIRequestScopesForbidden:           "Requested scopes are forbidden",
	ErrorAPIRequestBodyInvalid:               "Wrong input data",
	ErrorAPIRequestBodyParamsInvalid:         "Input data does not pass validation. Please specify valid params",
	ErrorAPIRequestBodyOldPasswordInvalid:    "Old password is invalid. Please check it again",
	ErrorAPIRequestBodyEmailInvalid:          "Specified email is invalid or empty",
	ErrorAPIRequestSignatureInvalid:          "Incorrect or empty request signature",
	ErrorAPIRequestAppIDInvalid:              "Incorrect or empty application ID",
	ErrorAPIRequestTokenInvalid:              "Incorrect or empty Bearer token",
	ErrorAPIAppInactive:                      "Requesting app is inactive",
	ErrorAPIAppRegistrationForbidden:         "Registration in this app is forbidden",
	ErrorAPIAppResetTokenNotCreated:          "Unable to create reset token",
	ErrorAPIAppAccessTokenNotCreated:         "Unable to create access token",
	ErrorAPIAppRefreshTokenNotCreated:        "Unable to create refresh token",
	ErrorAPIAppFederatedProviderNotSupported: "Federated provider is not supported",
	ErrorAPIAppFederatedProviderEmptyUserID:  "Federated provider returns empty user ID",
}

const (
	// ErrorAPIInternalServerError means that server got unknown error.
	ErrorAPIInternalServerError = "api.internal_server_error"
	// ErrorAPIUserUnableToCreate
	ErrorAPIUserUnableToCreate = "error.api.user.unable_to_create"
	// ErrorAPIUserNotFound
	ErrorAPIUserNotFound = "error.api.user.not_found"
	// ErrorAPIUsernameTaken
	ErrorAPIUsernameTaken = "error.api.username.taken"
	// ErrorAPIEmailTaken
	ErrorAPIEmailTaken = "error.api.email.taken"
	// ErrorAPIInviteTokenServerError
	ErrorAPIInviteTokenServerError = "error.api.invite_token.server_error"
	// ErrorAPIEmailNotSent
	ErrorAPIEmailNotSent = "error.api.email.not_sent"
	// ErrorAPIRequestPasswordWeak
	ErrorAPIRequestPasswordWeak = "error.api.request.password.weak"
	// ErrorAPIRequestScopesForbidden
	ErrorAPIRequestIncorrectEmailOrPassword = "error.api.request.incorrect_email_or_password"
	// ErrorAPIRequestScopesForbidden
	ErrorAPIRequestScopesForbidden = "error.api.request.scopes.forbidden"
	// ErrorAPIRequestBodyInvalid means that request data is corrupted.
	ErrorAPIRequestBodyInvalid = "error.api.request.body.invalid"
	// ErrorAPIRequestBodyParamsInvalid means that request data is corrupted.
	ErrorAPIRequestBodyParamsInvalid = "error.api.request.body.params.invalid"
	// ErrorAPIRequestBodyOldPasswordInvalid
	ErrorAPIRequestBodyOldPasswordInvalid = "error.api.request.body.oldpassword.invalid"
	// ErrorAPIRequestEmailInvalid means that email in request body is corrupted.
	ErrorAPIRequestBodyEmailInvalid = "error.api.request.body.email.invalid"
	// ErrorAPIRequestSignatureInvalid is a HMAC request signature error.
	ErrorAPIRequestSignatureInvalid = "error.api.request.signature.invalid"
	// ErrorAPIRequestAppIDInvalid means that application ID header value is invalid.
	ErrorAPIRequestAppIDInvalid = "error.api.request.app_id.invalid"
	// ErrorAPIRequestTokenInvalid means that the reqesting app is inactive.
	ErrorAPIRequestTokenInvalid = "error.api.request.token.invalid"
	// ErrorAPIAppInactive means that the token is invalid or empty.
	ErrorAPIAppInactive = "error.api.app.inactive"
	// ErrorAPIAppRegistrationForbidden means that registration is forbidden.
	ErrorAPIAppRegistrationForbidden = "error.api.app.registration_forbidden"
	// ErrorAPIAppResetTokenNotCreated means that registration is forbidden.
	ErrorAPIAppResetTokenNotCreated = "error.api.app.unable_to_create_reset_token"
	// ErrorAPIAppAccessTokenNotCreated means that registration is forbidden.
	ErrorAPIAppAccessTokenNotCreated = "error.api.app.unable_to_create_access_token"
	// ErrorAPIAppRefreshTokenNotCreated means that registration is forbidden.
	ErrorAPIAppRefreshTokenNotCreated = "error.api.app.unable_to_create_refresh_token"

	// ErrorAPIAppFederatedProviderNotSupported means that the federated ID provider is not supported.
	ErrorAPIAppFederatedProviderNotSupported = "api.app.federated.provider.not_supported"
	// ErrorAPIAppFederatedProviderEmptyUserID means that the federated provider returns empty user ID, maybe access token does not have required permissions.
	ErrorAPIAppFederatedProviderEmptyUserID = "api.app.federated.provider.empty_user_id"
)
