package admin

// Error is an http level error type.
type Error string

// Error is an implementation of std.Error interface.
func (e Error) Error() string { return string(e) }

const (
	// ErrorWrongInput is for corrupted request data.
	ErrorWrongInput = Error("Wrong input data")
	// ErrorRequestInvalidCookie is for invalid cookie.
	ErrorRequestInvalidCookie = Error("Invalid cookie")
	// ErrorInternalError is for internal errors.
	ErrorInternalError = Error("Internal error")
	// ErrorIncorrectLogin is for incorrect login and password.
	ErrorIncorrectLogin = Error("Incorrect login information")
	// ErrorNotAuthorized is for non-authorized access intents.
	ErrorNotAuthorized = Error("Not authorized")
	// ErrorAPIRequestBodyParamsInvalid means that request params are corrupted.
	ErrorAPIRequestBodyParamsInvalid = Error("Input data does not pass validation. Please specify valid params")
	// ErrorAPIInviteNotFound is when invite not found.
	ErrorAPIInviteNotFound = Error("Specified invite not found.")
	// ErrorAPIInviteUnableToInvalidate is when invite cannot be invalidated.
	ErrorAPIInviteUnableToInvalidate = Error("Unable to invalidate invite. Try again or contact support team")
	// ErrorAPIInviteTokenGenerate is when invite not found.
	ErrorAPIInviteTokenGenerate = Error("Unable to generate invite token")
	// ErrorAPISaveInvite is when invite not found.
	ErrorAPISaveInvite = Error("Unable to save invite")

	// ErrorAdminAccountIsNotSet is when not env variables for admin email or/and password are set
	ErrorAdminAccountIsNotSet = Error("Environment variabels for admin account (email and password) is not set")
	// ErrorAdminAccountNoEmailAndPassword when no admin email or/and password are set
	ErrorAdminAccountNoEmailAndPassword = Error("Admin email and/or password are empty and could not be used for security reason")
)
