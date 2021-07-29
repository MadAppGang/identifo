package api

// MessageID is an error message ID.
type MessageID string

// GetMessage returns message by its ID.
func GetMessage(id MessageID) string {
	return messages[id]
}

var messages = map[MessageID]string{
	ErrorAPIInternalServerError:              "Internal server error",
	ErrorAPIUserUnableToCreate:               "Unable to create use. Try again or contact support team",
	ErrorAPIVerificationCodeInvalid:          "Sorry, the code you entered is invalid or has expired. Please get a new one.",
	ErrorAPIUserNotFound:                     "Specified user not found",
	ErrorAPIUsernameTaken:                    "Username is taken. Try to choose another one",
	ErrorAPIEmailTaken:                       "Email is taken. Try to choose another one",
	ErrorAPIInviteTokenServerError:           "Unable to create invite token. Try again or contact support team",
	ErrorAPIInviteUnableToInvalidate:         "Unable to invalidate invite. Try again or contact support team",
	ErrorAPIInviteUnableToSave:               "Unable to save invite. Try again or contact support team",
	ErrorAPIInviteUnableToGet:                "Unable to get invites. Try again or contact support team",
	ErrorAPIEmailNotSent:                     "Unable to send email. Try again or contact support team",
	ErrorAPIRequestPasswordWeak:              "Password is not strong enough",
	ErrorAPIRequestIncorrectLoginOrPassword:  "Incorrect email or password",
	ErrorAPIRequestScopesForbidden:           "Requested scopes are forbidden",
	ErrorAPIRequestBodyInvalid:               "Wrong input data",
	ErrorAPIRequestBodyParamsInvalid:         "Input data does not pass validation. Please specify valid params",
	ErrorAPIRequestBodyOldPasswordInvalid:    "Old password is invalid. Please check it again",
	ErrorAPIRequestBodyEmailInvalid:          "Specified email is invalid or empty",
	ErrorAPIRequestSignatureInvalid:          "Incorrect or empty request signature",
	ErrorAPIRequestAppIDInvalid:              "Incorrect or empty application ID",
	ErrorAPIRequestTokenInvalid:              "Incorrect or empty Bearer token",
	ErrorAPIRequestTFACodeEmpty:              "Empty two-factor authentication code",
	ErrorAPIRequestTFACodeInvalid:            "Invalid two-factor authentication code",
	ErrorAPIRequestTFAAlreadyEnabled:         "Two-factor authentication already enabled",
	ErrorAPIRequestPleaseEnableTFA:           "Please enable two-factor authenticaton",
	ErrorAPIRequestPleaseDisableTFA:          "Please disable two-factor authenticaton",
	ErrorAPIRequestMandatoryTFA:              "Two-factor authentication is mandatory for this app",
	ErrorAPIRequestDisabledTFA:               "Two-factor authentication is disabled for this app",
	ErrorAPIRequestPleaseSetPhoneForTFA:      "Please specify your phone number to be able to receive one-time passwords",
	ErrorAPIRequestPleaseSetEmailForTFA:      "Please specify your email address to be able to receive one-time passwords",
	ErrorAPIAppInactive:                      "Requesting app is inactive",
	ErrorAPIAppRegistrationForbidden:         "Registration in this app is forbidden",
	ErrorAPIAppResetTokenNotCreated:          "Unable to create reset token",
	ErrorAPIAppAccessTokenNotCreated:         "Unable to create access token",
	ErrorAPIAppRefreshTokenNotCreated:        "Unable to create refresh token",
	ErrorAPIAppCannotExtractTokenSubject:     "Unable to extract Subject claim from token",
	ErrorAPIAppCannotInitAuthorizer:          "Unable to init internal authorizer",
	ErrorAPIAppFederatedProviderNotSupported: "Federated provider is not supported",
	ErrorAPIAppLoginWithUsernameNotSupported: "Login with username is not supported by app",
	ErrorAPIAppPhoneLoginNotSupported:        "Login with phone number is not supported by app",
	ErrorAPIAppAccessDenied:                  "Access denied",
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
	// ErrorAPIInviteUnableToInvalidate is when invite cannot be invalidated.
	ErrorAPIInviteUnableToInvalidate = "error.api.invite.unable_to_invalidate"
	// ErrorAPIInviteUnableToSave is when invite cannot be saved.
	ErrorAPIInviteUnableToSave = "error.api.invite.unable_to_save"
	// ErrorAPIInviteUnableToGet is when invites cannot be fetched.
	ErrorAPIInviteUnableToGet = "errors.api.invite.unable_to_get"
	// ErrorAPIEmailNotSent means that email had not been sent.
	ErrorAPIEmailNotSent = "error.api.email.not_sent"

	// ErrorAPIRequestPasswordWeak means that password didn't pass strength validation.
	ErrorAPIRequestPasswordWeak = "error.api.request.password.weak"
	// ErrorAPIRequestIncorrectLoginOrPassword is for incorrect login or password.
	ErrorAPIRequestIncorrectLoginOrPassword = "error.api.request.incorrect_login_or_password"
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

	// ErrorAPIRequestTFACodeEmpty means that the 2FA code is empty.
	ErrorAPIRequestTFACodeEmpty = "error.api.request.2fa_code.empty"
	// ErrorAPIRequestTFACodeInvalid means that the 2FA code is invalid.
	ErrorAPIRequestTFACodeInvalid = "error.api.request.2fa_code.invalid"
	// ErrorAPIRequestTFAAlreadyEnabled means that 2FA is already enabled for the user.
	ErrorAPIRequestTFAAlreadyEnabled = "error.api.request.2fa.already_enabled"
	// ErrorAPIRequestPleaseEnableTFA means that user must request TFA and obtain TFA secret to be able to use the app.
	ErrorAPIRequestPleaseEnableTFA = "error.api.request.2fa.please_enable"
	// ErrorAPIRequestPleaseDisableTFA means that user must disable TFA to be able to use the app.
	ErrorAPIRequestPleaseDisableTFA = "error.api.request.2fa.please_disable"
	// ErrorAPIRequestMandatoryTFA means that user cannot disable TFA for the app.
	ErrorAPIRequestMandatoryTFA = "error.api.request.2fa.mandatory"
	// ErrorAPIRequestDisabledTFA means that app does not support TFA.
	ErrorAPIRequestDisabledTFA = "error.api.request.2fa.disabled"
	// ErrorAPIRequestPleaseSetPhoneForTFA means that user must set up their phone number to be able to receive OTPs in SMS.
	ErrorAPIRequestPleaseSetPhoneForTFA = "error.api.request.2fa.set_phone"
	// ErrorAPIRequestPleaseSetEmailForTFA means that user must set up their email address to be able to receive OTPs on the email.
	ErrorAPIRequestPleaseSetEmailForTFA = "error.api.request.2fa.set_email"
	// ErrorAPIRequestUnableToSendOTP means that there is error sending the otp code while login to user
	ErrorAPIRequestUnableToSendOTP = "error.api.request.2fa.unable to send OTP code to email or sms"

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
	// ErrorAPIAppCannotExtractTokenSubject is when we cannot extract token "sub".
	ErrorAPIAppCannotExtractTokenSubject = "error.api.request.token.sub"
	// ErrorAPIAppCannotInitAuthorizer is when we cannot init internal authorizer.
	ErrorAPIAppCannotInitAuthorizer = "error.api.request.authorizer.internal.init"

	// ErrorAPIAppFederatedProviderNotSupported means that the federated ID provider is not supported.
	ErrorAPIAppFederatedProviderNotSupported = "api.app.federated.provider.not_supported"
	// ErrorAPIAppFederatedProviderEmptyUserID means that the federated provider returns empty user ID, maybe access token does not have required permissions.
	ErrorAPIAppFederatedEmptyRedirect = "api.app.federated.provider.empty_redirect"
	ErrorAPIAppFederatedEmptyProvider = "api.app.federated.provider.empty"
	ErrorAPIAppFederatedCantComplete  = "api.app.federated.provider.cant_complete"
	// ErrorAPIAppFederatedProviderEmptyAppleInfo means that application does not have clientID and clientSecret needed for Sign In with Apple.
	// ErrorAPIAppFederatedProviderEmptyAppleInfo = "api.app.federated.provider.empty_apple_info"

	// ErrorAPIAppFederatedLoginNotSupported means that the app does not support federated login.
	// ErrorAPIAppFederatedLoginNotSupported = "api.app.federated.login.not_supported"
	// ErrorAPIAppLoginWithUsernameNotSupported means that the app does not support login by username.
	ErrorAPIAppLoginWithUsernameNotSupported = "api.app.username.login.not_supported"
	// ErrorAPIAppPhoneLoginNotSupported means that the app does not support login by phone number.
	ErrorAPIAppPhoneLoginNotSupported = "api.app.phone.login.not_supported"
)
