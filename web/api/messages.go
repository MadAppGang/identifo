package api

// MessageID is an error message ID.
type MessageID string

// GetMessage returns message by its ID.
func GetMessage(id MessageID) string {
	return messages[id]
}

var messages = map[MessageID]string{
	ErrorAPIInternalServerError:                "Internal server error",
	ErrorAPIUserUnableToCreate:                 "Unable to create use. Try again or contact support team",
	ErrorAPIVerificationCodeInvalid:            "Sorry, the code you entered is invalid or has expired. Please get a new one.",
	ErrorAPIUserNotFound:                       "Specified user not found",
	ErrorAPIUsernameTaken:                      "Username is taken. Try to choose another one",
	ErrorAPIEmailTaken:                         "Email is taken. Try to choose another one",
	ErrorAPIInviteTokenServerError:             "Unable to create invite token. Try again or contact support team",
	ErrorAPIEmailNotSent:                       "Unable to send email. Try again or contact support team",
	ErrorAPIRequestPasswordWeak:                "Password is not strong enough",
	ErrorAPIRequestIncorrectEmailOrPassword:    "Incorrect email or password",
	ErrorAPIRequestScopesForbidden:             "Requested scopes are forbidden",
	ErrorAPIRequestBodyInvalid:                 "Wrong input data",
	ErrorAPIRequestBodyParamsInvalid:           "Input data does not pass validation. Please specify valid params",
	ErrorAPIRequestBodyOldPasswordInvalid:      "Old password is invalid. Please check it again",
	ErrorAPIRequestBodyEmailInvalid:            "Specified email is invalid or empty",
	ErrorAPIRequestSignatureInvalid:            "Incorrect or empty request signature",
	ErrorAPIRequestAppIDInvalid:                "Incorrect or empty application ID",
	ErrorAPIRequestTokenInvalid:                "Incorrect or empty Bearer token",
	ErrorAPIAppInactive:                        "Requesting app is inactive",
	ErrorAPIAppRegistrationForbidden:           "Registration in this app is forbidden",
	ErrorAPIAppResetTokenNotCreated:            "Unable to create reset token",
	ErrorAPIAppAccessTokenNotCreated:           "Unable to create access token",
	ErrorAPIAppRefreshTokenNotCreated:          "Unable to create refresh token",
	ErrorAPIAppCannotExtractTokenSubject:       "Unable to extract Subject claim from token",
	ErrorAPIAppCannotInitAuthorizer:            "Unable to init internal authorizer",
	ErrorAPIAppFederatedProviderNotSupported:   "Federated provider is not supported",
	ErrorAPIAppFederatedProviderEmptyUserID:    "Federated provider returns empty user ID",
	ErrorAPIAppFederatedProviderEmptyAppleInfo: "Application does not have Apple info",
	ErrorAPIAppFederatedLoginNotSupported:      "Login with federated identity provider is not supported by app",
	ErrorAPIAppLoginWithUsernameNotSupported:   "Login with username is not supported by app",
	ErrorAPIAppPhoneLoginNotSupported:          "Login with phone number is not supported by app",
	ErrorAPIAppAccessDenied:                    "Access denied",
}

const (
	// ErrorAPIInternalServerError means that server got unknown error.
	ErrorAPIInternalServerError = "api.internal_server_error"
	// ErrorAPIAppAccessDenied is when access is denied.
	ErrorAPIAppAccessDenied = "api.app.access_denied"
	// ErrorAPIUserUnableToCreate is when user cannot create the resource.
	ErrorAPIUserUnableToCreate = "error.api.user.unable_to_create"
	// ErrorAPIVerificationCodeInvalid stands for invalid verification code.
	ErrorAPIVerificationCodeInvalid = "error.api.verification_code.invalid"
	// ErrorAPIUserNotFound is when user not found.
	ErrorAPIUserNotFound = "error.api.user.not_found"
	// ErrorAPIUsernameTaken is when username is already taken.
	ErrorAPIUsernameTaken = "error.api.username.taken"
	// ErrorAPIEmailTaken is when email is already taken.
	ErrorAPIEmailTaken = "error.api.email.taken"
	// ErrorAPIInviteTokenServerError is for invite token creation issues.
	ErrorAPIInviteTokenServerError = "error.api.invite_token.server_error"
	// ErrorAPIEmailNotSent means that email had not been sent.
	ErrorAPIEmailNotSent = "error.api.email.not_sent"

	// ErrorAPIRequestPasswordWeak means that password didn't pass strength validation.
	ErrorAPIRequestPasswordWeak = "error.api.request.password.weak"
	// ErrorAPIRequestIncorrectEmailOrPassword is for incorrect email or password.
	ErrorAPIRequestIncorrectEmailOrPassword = "error.api.request.incorrect_email_or_password"
	// ErrorAPIRequestScopesForbidden is for forbidden request scopes.
	ErrorAPIRequestScopesForbidden = "error.api.request.scopes.forbidden"
	// ErrorAPIRequestBodyInvalid means that request body is corrupted.
	ErrorAPIRequestBodyInvalid = "error.api.request.body.invalid"
	// ErrorAPIRequestBodyParamsInvalid means that request params are corrupted.
	ErrorAPIRequestBodyParamsInvalid = "error.api.request.body.params.invalid"
	// ErrorAPIRequestBodyOldPasswordInvalid is for invalid old password.
	ErrorAPIRequestBodyOldPasswordInvalid = "error.api.request.body.oldpassword.invalid"
	// ErrorAPIRequestBodyEmailInvalid means that email in request body is corrupted.
	ErrorAPIRequestBodyEmailInvalid = "error.api.request.body.email.invalid"
	// ErrorAPIRequestSignatureInvalid is a HMAC request signature error.
	ErrorAPIRequestSignatureInvalid = "error.api.request.signature.invalid"
	// ErrorAPIRequestAppIDInvalid means that application ID header value is invalid.
	ErrorAPIRequestAppIDInvalid = "error.api.request.app_id.invalid"
	// ErrorAPIRequestTokenInvalid means that the token is invalid or empty.
	ErrorAPIRequestTokenInvalid = "error.api.request.token.invalid"
	// ErrorAPIAppCannotExtractTokenSubject is when we cannot extract token "sub".
	ErrorAPIAppCannotExtractTokenSubject = "error.api.request.token.sub"
	// ErrorAPIAppCannotInitAuthorizer is when we cannot init internal authorizer.
	ErrorAPIAppCannotInitAuthorizer = "error.api.request.authorizer.internal.init"

	// ErrorAPIAppInactive means that the reqesting app is inactive.
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
	// ErrorAPIAppFederatedProviderEmptyAppleInfo means that application does not have clientID and clientSecret needed for Sign In with Apple.
	ErrorAPIAppFederatedProviderEmptyAppleInfo = "api.app.federated.provider.empty_apple_info"

	// ErrorAPIAppFederatedLoginNotSupported means that the app does not support federated login.
	ErrorAPIAppFederatedLoginNotSupported = "api.app.federated.login.not_supported"
	// ErrorAPIAppLoginWithUsernameNotSupported means that the app does not support login by username.
	ErrorAPIAppLoginWithUsernameNotSupported = "api.app.username.login.not_supported"
	// ErrorAPIAppPhoneLoginNotSupported means that the app does not support login by phone number.
	ErrorAPIAppPhoneLoginNotSupported = "api.app.phone.login.not_supported"
)
