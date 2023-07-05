// Code generated by "./gen.go"; DO NOT EDIT.

package l

type LocalizedString string

func (l LocalizedString) Error() string {
	return string(l)
}

const (

	//===========================================================================
	//  test messages
	//===========================================================================
	// Test -> I am test string
	Test LocalizedString = "test"
	// KeyWithSpaces -> key with spaces passed test
	KeyWithSpaces LocalizedString = "key with spaces"

	//===========================================================================
	//  server settings messages
	//===========================================================================
	// ErrorServerUpdateSettingsNoChange -> no settings has been changed, skipping the update
	ErrorServerUpdateSettingsNoChange LocalizedString = "error.server.update.settings.no.change"
	// ErrorServerUpdateSettingsValidationError -> settings validation failed with error: %v
	ErrorServerUpdateSettingsValidationError LocalizedString = "error.server.update.settings.validation.error"
	// ErrorServerUpdateSettingsAPPlyError -> error saving new config: %v
	ErrorServerUpdateSettingsAPPlyError LocalizedString = "error.server.update.settings.apply.error"

	//===========================================================================
	//  api messages
	//===========================================================================
	// APIInternalServerError -> Internal server error.
	APIInternalServerError LocalizedString = "api.internal_server_error"
	// APIInternalServerErrorWithError -> Internal server error; %v.
	APIInternalServerErrorWithError LocalizedString = "api.internal_server_error_with_error"
	// APIAccessDenied -> Access denied.
	APIAccessDenied LocalizedString = "api.access_denied"
	// ErrorAPIUserUnableToCreate -> Unable to create use, please try again or contact support team.
	ErrorAPIUserUnableToCreate LocalizedString = "error.api.user.unable_to_create"
	// ErrorAPIVerificationCodeInvalid -> Sorry, the code you entered is invalid or has expired. Please get a new one.
	ErrorAPIVerificationCodeInvalid LocalizedString = "error.api.verification_code.invalid"
	// ErrorAPIUserNotFound -> User not found.
	ErrorAPIUserNotFound LocalizedString = "error.api.user.not_found"
	// ErrorAPIUserNotFoundError -> User not found with error; %v.
	ErrorAPIUserNotFoundError LocalizedString = "error.api.user.not_found.error"
	// ErrorAPIUsernameTaken -> Username is taken. Try to choose another one.
	ErrorAPIUsernameTaken LocalizedString = "error.api.username.taken"
	// ErrorAPIUsernamePhoneEmailTaken -> Username, email or/and phone is taken. Try to choose another one.
	ErrorAPIUsernamePhoneEmailTaken LocalizedString = "error.api.username_phone_email.taken"
	// ErrorAPIEmailTaken -> Email is taken. Try to choose another one.
	ErrorAPIEmailTaken LocalizedString = "error.api.email.taken"
	// ErrorAPIPhoneTaken -> Phone is taken. Try to choose another one.
	ErrorAPIPhoneTaken LocalizedString = "error.api.phone.taken"
	// ErrorAPIInviteUnableToInvalidateError -> Bad invite token, please try again or contact support team. Token error: %v.
	ErrorAPIInviteUnableToInvalidateError LocalizedString = "error.api.invite.unable_to_invalidate.error"
	// ErrorAPIInviteUnableToSave -> Unable to save invite. Try again or contact support team.
	ErrorAPIInviteUnableToSave LocalizedString = "error.api.invite.unable_to_save"
	// ErrorAPIInviteUnableToGet -> Unable to get invites. Try again or contact support team.
	ErrorAPIInviteUnableToGet LocalizedString = "error.api.invite.unable_to_get"
	// ErrorAPIEmailNotSent -> Unable to send email. Try again or contact support team.
	ErrorAPIEmailNotSent LocalizedString = "error.api.email.not_sent"
	// ErrorAPIRequestPasswordWeak -> Password is not strong enough: %v
	ErrorAPIRequestPasswordWeak LocalizedString = "error.api.request.password.weak"
	// ErrorAPIRequestIncorrectLoginOrPassword -> Invalid Username or Password!
	ErrorAPIRequestIncorrectLoginOrPassword LocalizedString = "error.api.request.incorrect_login_or_password"
	// ErrorAPIRequestScopesForbidden -> Requested scopes are forbidden.
	ErrorAPIRequestScopesForbidden LocalizedString = "error.api.request.scopes.forbidden"
	// ErrorAPIRequestBodyInvalidError -> Error reading request body data: %v.
	ErrorAPIRequestBodyInvalidError LocalizedString = "error.api.request.body.invalid.error"
	// ErrorAPIRequestBodyOldpasswordInvalid -> Old password is invalid. Please check it and try again.
	ErrorAPIRequestBodyOldpasswordInvalid LocalizedString = "error.api.request.body.oldpassword.invalid"
	// ErrorAPIRequestBodyEmailInvalid -> Specified email is invalid or empty.
	ErrorAPIRequestBodyEmailInvalid LocalizedString = "error.api.request.body.email.invalid"
	// ErrorAPIRequestSignatureInvalid -> Incorrect or empty request signature.
	ErrorAPIRequestSignatureInvalid LocalizedString = "error.api.request.signature.invalid"
	// ErrorAPIRequestSignatureValidationError -> Incorrect request signature: %v.
	ErrorAPIRequestSignatureValidationError LocalizedString = "error.api.request.signature.validation.error"
	// ErrorAPIRequestAPPIDInvalid -> Incorrect or empty application ID.
	ErrorAPIRequestAPPIDInvalid LocalizedString = "error.api.request.app_id.invalid"
	// ErrorAPIRequestCallbackurlInvalid -> Please add callbackURL in application settings.
	ErrorAPIRequestCallbackurlInvalid LocalizedString = "error.api.request.callbackurl.invalid"
	// ErrorAPISessionNotFound -> Unable find a matching session for this request: %s.
	ErrorAPISessionNotFound LocalizedString = "error.api.session.not.found"
	// ErrorAPILoginError -> Login error: %v.
	ErrorAPILoginError LocalizedString = "error.api.login.error"
	// ErrorAPILoginCodeInvalid -> The code you entered is incorrect. Please check it and try again.
	ErrorAPILoginCodeInvalid LocalizedString = "error.api.login.code.invalid"
	// ErrorAPILoginAnonymousForbidden -> Anonymous login is forbidden for this app.
	ErrorAPILoginAnonymousForbidden LocalizedString = "error.api.login.anonymous.forbidden"
	// ErrorAPIInviteEmailMismatch -> Invite email and user email are not equal.
	ErrorAPIInviteEmailMismatch LocalizedString = "error.api.invite.email.mismatch"
	// ErrorAPIInviteRoleMissing -> No role in invite token found.
	ErrorAPIInviteRoleMissing LocalizedString = "error.api.invite.role.missing"
	// ErrorAPIDataError -> Error getting data from request: %s
	ErrorAPIDataError LocalizedString = "error.api.data.error"

	//===========================================================================
	// general api errors
	//===========================================================================
	// ErrorAPIJsonParseError -> Error parsing request JSON: %v.
	ErrorAPIJsonParseError LocalizedString = "error.api.json.parse.error"

	//===========================================================================
	//  2FA errors
	//===========================================================================
	// ErrorAPIRequest2FACodeEmpty -> Empty two-factor authentication code.
	ErrorAPIRequest2FACodeEmpty LocalizedString = "error.api.request.2fa_code.empty"
	// ErrorAPIRequest2FACodeInvalid -> Invalid two-factor authentication code.
	ErrorAPIRequest2FACodeInvalid LocalizedString = "error.api.request.2fa_code.invalid"
	// ErrorAPIRequest2FAAlreadyEnabled -> Two-factor authentication already enabled.
	ErrorAPIRequest2FAAlreadyEnabled LocalizedString = "error.api.request.2fa.already_enabled"
	// ErrorAPIRequest2FAPleaseEnable -> Please enable two-factor authentication.
	ErrorAPIRequest2FAPleaseEnable LocalizedString = "error.api.request.2fa.please_enable"
	// ErrorAPIRequest2FAPleaseDisable -> Please disable two-factor authentication.
	ErrorAPIRequest2FAPleaseDisable LocalizedString = "error.api.request.2fa.please_disable"
	// ErrorAPIRequest2FAMandatory -> Two-factor authentication is required for this app.
	ErrorAPIRequest2FAMandatory LocalizedString = "error.api.request.2fa.mandatory"
	// ErrorAPIRequest2FADisabled -> Two-factor authentication is disabled for this app.
	ErrorAPIRequest2FADisabled LocalizedString = "error.api.request.2fa.disabled"
	// ErrorAPIRequest2FASetPhone -> Please specify your phone number to be able to receive one-time passwords.
	ErrorAPIRequest2FASetPhone LocalizedString = "error.api.request.2fa.set_phone"
	// ErrorAPIRequest2FASetEmail -> Please specify your email address to be able to receive one-time passwords.
	ErrorAPIRequest2FASetEmail LocalizedString = "error.api.request.2fa.set_email"
	// ErrorAPIRequestEnable2FAEmptyPhoneAndEmail -> Phone and email are empty.
	ErrorAPIRequestEnable2FAEmptyPhoneAndEmail LocalizedString = "error.api.request.enable_2fa.empty_phone_and_email"
	// ErrorAPIRequest2FAUnableToSendOtpError -> Error sending OTP code with SMS or Email with error: %v.
	ErrorAPIRequest2FAUnableToSendOtpError LocalizedString = "error.api.request.2fa.unable_to_send_OTP.error"
	// ErrorAPIRequest2FAUnableToGenerateQrError -> Unable to create QR code with error: %v.
	ErrorAPIRequest2FAUnableToGenerateQrError LocalizedString = "error.api.request.2fa.unable_to_generate_QR.error"
	// ErrorAPIRequest2FAUnknownType -> Unknown TFA type: %s.
	ErrorAPIRequest2FAUnknownType LocalizedString = "error.api.request.2fa.unknown_type"
	// Error2FAResendTimeout -> Please wait before new code resend.
	Error2FAResendTimeout LocalizedString = "error.2fa.resend.timeout"
	// Error2FAVerifyFailError -> OTP code is invalid: %v
	Error2FAVerifyFailError LocalizedString = "error.2fa.verify.fail.error"

	//===========================================================================
	//  Token errors
	//===========================================================================
	// ErrorAPIRequestTokenSub -> Unable to extract Subject claim from the token.
	ErrorAPIRequestTokenSub LocalizedString = "error.api.request.token.sub"
	// ErrorAPIRequestTokenSubError -> Unable to extract Subject claim from the token with error: %v.
	ErrorAPIRequestTokenSubError LocalizedString = "error.api.request.token.sub.error"
	// ErrorAPIRequestTokenInvalid -> Incorrect or empty Bearer token.
	ErrorAPIRequestTokenInvalid LocalizedString = "error.api.request.token.invalid"
	// ErrorAPIContextNoToken -> Error getting token from context.
	ErrorAPIContextNoToken LocalizedString = "error.api.context.no_token"
	// ErrorAPITokenParseError -> Error parsing access token: %v.
	ErrorAPITokenParseError LocalizedString = "error.api.token.parse.error"
	// ErrorTokenInviteCreateError -> Unable to create invite token with error: %v.
	ErrorTokenInviteCreateError LocalizedString = "error.token.invite.create.error"
	// ErrorTokenUnableToCreateResetTokenError -> Error creating reset token: %v.
	ErrorTokenUnableToCreateResetTokenError LocalizedString = "error.token.unable_to_create_reset_token.error"
	// ErrorTokenUnableToCreateAccessTokenError -> Error creating access token: %v.
	ErrorTokenUnableToCreateAccessTokenError LocalizedString = "error.token.unable_to_create_access_token.error"
	// ErrorTokenUnableToCreateRefreshTokenError -> Error creating refresh token: %v.
	ErrorTokenUnableToCreateRefreshTokenError LocalizedString = "error.token.unable_to_create_refresh_token.error"
	// ErrorTokenUnableToCreateInviteTokenError -> Error creating invite token: %v.
	ErrorTokenUnableToCreateInviteTokenError LocalizedString = "error.token.unable_to_create_invite_token.error"
	// ErrorTokenRefreshAccessToken -> Error getting new access token with refresh token: %v.
	ErrorTokenRefreshAccessToken LocalizedString = "error.token.refresh_access_token"
	// ErrorTokenRefreshEmpty -> Error getting old refresh token from context to replace it.
	ErrorTokenRefreshEmpty LocalizedString = "error.token.refresh.empty"
	// ErrorOtpExpired -> OTP token expired, please get the new one and try again.
	ErrorOtpExpired LocalizedString = "error.otp.expired"
	// ErrorTokenInvalidError -> Invalid token. Validation error: %v.
	ErrorTokenInvalidError LocalizedString = "error.token.invalid.error"
	// ErrorTokenBlocked -> The token is blocked and not valid any more.
	ErrorTokenBlocked LocalizedString = "error.token.blocked"

	//===========================================================================
	//  App errors
	//===========================================================================
	// ErrorAPIAPPInactive -> The app is inactive.
	ErrorAPIAPPInactive LocalizedString = "error.api.app.inactive"
	// ErrorAPIAPPRegistrationForbidden -> Registration in this app is forbidden.
	ErrorAPIAPPRegistrationForbidden LocalizedString = "error.api.app.registration_forbidden"
	// ErrorAPIRequestAuthorizerInternalInit -> Error creating authz service.
	ErrorAPIRequestAuthorizerInternalInit LocalizedString = "error.api.request.authorizer.internal.init"
	// ErrorAPIAPPUnableToTokenPayloadForAPPError -> Error getting token payload for the app %s with error: %v.
	ErrorAPIAPPUnableToTokenPayloadForAPPError LocalizedString = "error.api.app.unable_to_token_payload_for_app.error"
	// ErrorAPIAPPNoAPPInContext -> Missing app data in context.
	ErrorAPIAPPNoAPPInContext LocalizedString = "error.api.app.no_app_in_context"
	// ErrorAPPRegisterUrlError -> Invalid register URL (%s) for app (%s): %v.
	ErrorAPPRegisterUrlError LocalizedString = "error.app.register_url.error"
	// ErrorAPPLoginNoScope -> User has no required scopes by this app.
	ErrorAPPLoginNoScope LocalizedString = "error.app.login.no_scope"
	// ErrorAPPResetUrlError -> Invalid reset password URL (%s) for app (%s): %v.
	ErrorAPPResetUrlError LocalizedString = "error.app.reset_url.error"

	//===========================================================================
	//  Login
	//===========================================================================
	// ErrorUpdateUserImmutableField -> Field %s could not be updated.
	ErrorUpdateUserImmutableField LocalizedString = "error.update.user.immutable.field"
	// ErrorUpdateUserFieldValueAlreadyOccupied -> %s %s already exists, choose other one
	ErrorUpdateUserFieldValueAlreadyOccupied LocalizedString = "error.update.user.field.value.already.occupied"
	// ErrorInvalidPhone -> invalid phone number
	ErrorInvalidPhone LocalizedString = "error.invalid.phone"
	// ErrorEmailEmpty -> email could not be empty
	ErrorEmailEmpty LocalizedString = "error.email.empty"
	// ErrorPhoneEmpty -> phone could not be empty
	ErrorPhoneEmpty LocalizedString = "error.phone.empty"
	// ErrorUsernameEmpty -> username could not be empty
	ErrorUsernameEmpty LocalizedString = "error.username.empty"

	//===========================================================================
	//  Federated login
	//===========================================================================
	// APIAPPFederatedProviderNotSupported -> Federated provider is not supported.
	APIAPPFederatedProviderNotSupported LocalizedString = "api.app.federated.provider.not_supported"
	// APIAPPFederatedProviderEmptyRedirect -> Empty redirect URL.
	APIAPPFederatedProviderEmptyRedirect LocalizedString = "api.app.federated.provider.empty_redirect"
	// APIAPPFederatedProviderEmpty -> Empty federated login provider.
	APIAPPFederatedProviderEmpty LocalizedString = "api.app.federated.provider.empty"
	// APIAPPFederatedProviderCantCompleteError -> Unable to complete federated login: %v.
	APIAPPFederatedProviderCantCompleteError LocalizedString = "api.app.federated.provider.cant_complete.error"
	// APIFederatedCreateAuthUrlError -> Unable to create auth URL with error: %v.
	APIFederatedCreateAuthUrlError LocalizedString = "api.federated.create_auth_url.error"
	// APIAPPUsernameLoginNotSupported -> Login with username is not supported by app.
	APIAPPUsernameLoginNotSupported LocalizedString = "api.app.username.login.not_supported"
	// APIAPPPhoneLoginNotSupported -> Login with phone number is not supported by app.
	APIAPPPhoneLoginNotSupported LocalizedString = "api.app.phone.login.not_supported"
	// ErrorAPIUnableToInitializeIDentifo -> Unable to initialize NativeLogin.
	ErrorAPIUnableToInitializeIDentifo LocalizedString = "error.api.unable_to_initialize_identifo"
	// ErrorFederatedUnmarshalSessionError -> Error getting federated login session: %v.
	ErrorFederatedUnmarshalSessionError LocalizedString = "error.federated.unmarshal.session.error"
	// ErrorFederatedSessionAPPIDMismatch -> Session app id(%s) and request app id(%s) mismatch.
	ErrorFederatedSessionAPPIDMismatch LocalizedString = "error.federated.session_app_id_mismatch"
	// ErrorFederatedAccessDeniedError -> You are not allowed to login with error: %v.
	ErrorFederatedAccessDeniedError LocalizedString = "error.federated.access_denied.error"
	// ErrorFederatedLoginError -> Federated login error: %v.
	ErrorFederatedLoginError LocalizedString = "error.federated.login.error"
	// ErrorFederatedCodeError -> No code returned for federated login
	ErrorFederatedCodeError LocalizedString = "error.federated.code.error"
	// ErrorFederatedStateError -> State mismatch code returned for federated login
	ErrorFederatedStateError LocalizedString = "error.federated.state.error"
	// ErrorFederatedExchangeError -> Federated exchange error: %v.
	ErrorFederatedExchangeError LocalizedString = "error.federated.exchange.error"
	// ErrorFederatedIDtokenMissing -> No id_token returned for federated login
	ErrorFederatedIDtokenMissing LocalizedString = "error.federated.idtoken.missing"
	// ErrorFederatedIDtokenInvalid -> Invalid id_token returned for federated login: %v
	ErrorFederatedIDtokenInvalid LocalizedString = "error.federated.idtoken.invalid"
	// ErrorFederatedClaimsError -> Invalid claims error: %v
	ErrorFederatedClaimsError LocalizedString = "error.federated.claims.error"
	// ErrorFederatedOidcProviderError -> Failed to init OIDC provider: %v
	ErrorFederatedOidcProviderError LocalizedString = "error.federated.oidc.provider.error"
	// ErrorFederatedOidcDisabled -> Federated OIDC login disabled
	ErrorFederatedOidcDisabled LocalizedString = "error.federated.oidc.disabled"

	//===========================================================================
	//  Storages
	//===========================================================================
	// ErrorStorageConnectionError -> Unable to connect to storage with error: %v.
	ErrorStorageConnectionError LocalizedString = "error.storage.connection.error"
	// ErrorStorageUnsupportedType -> Error connecting to server, unsupported connection type.
	ErrorStorageUnsupportedType LocalizedString = "error.storage.unsupported.type"
	// ErrorStorageUpdateUserError -> Unable to update user with id %s with error: %v
	ErrorStorageUpdateUserError LocalizedString = "error.storage.update_user.error"
	// ErrorStorageFindUserEmailError -> Unable to find user with email %s with error: %v
	ErrorStorageFindUserEmailError LocalizedString = "error.storage.find.user.email.error"
	// ErrorStorageFindUserIDError -> Unable to find user with id %s with error: %v
	ErrorStorageFindUserIDError LocalizedString = "error.storage.find.user.id.error"
	// ErrorStorageFindUserPhoneError -> Unable to find user with phone %s with error: %v
	ErrorStorageFindUserPhoneError LocalizedString = "error.storage.find.user.phone.error"
	// ErrorStorageFindUserEmailPhoneUsernameError -> Unable to find user with error: %v.
	ErrorStorageFindUserEmailPhoneUsernameError LocalizedString = "error.storage.find.user.email_phone_username.error"
	// ErrorStorageResetPasswordUserError -> Error saving new password for user(id:%s): %v.
	ErrorStorageResetPasswordUserError LocalizedString = "error.storage.reset_password.user.error"
	// ErrorStorageAPPFindByIDError -> Unable to find app with id %s with error: %v.
	ErrorStorageAPPFindByIDError LocalizedString = "error.storage.app.find.by_id.error"
	// ErrorStorageUserFederatedCreateError -> Error creating federated user: %v.
	ErrorStorageUserFederatedCreateError LocalizedString = "error.storage.user.federated.create.error"
	// ErrorStorageUserCreateError -> Error creating user: %v.
	ErrorStorageUserCreateError LocalizedString = "error.storage.user.create.error"
	// ErrorStorageInviteFindEmailError -> Error getting invite by email: %v.
	ErrorStorageInviteFindEmailError LocalizedString = "error.storage.invite.find.email.error"
	// ErrorStorageInviteFindIDError -> Error getting invite by id: %v.
	ErrorStorageInviteFindIDError LocalizedString = "error.storage.invite.find.id.error"
	// ErrorStorageInviteArchiveEmailError -> Error archiving old invited by email: %v.
	ErrorStorageInviteArchiveEmailError LocalizedString = "error.storage.invite.archive.email.error"
	// ErrorStorageInviteSaveError -> Error saving invite token: %v.
	ErrorStorageInviteSaveError LocalizedString = "error.storage.invite.save.error"
	// ErrorStorageVerificationCreateError -> Error creating phone verification code: %v.
	ErrorStorageVerificationCreateError LocalizedString = "error.storage.verification.create.error"
	// ErrorStorageVerificationFindError -> Error getting verification code from storage: %v.
	ErrorStorageVerificationFindError LocalizedString = "error.storage.verification.find.error"
	// ErrorStorageRequestError -> error making request to db: %v
	ErrorStorageRequestError LocalizedString = "error.storage.request.error"

	//===========================================================================
	//  Services
	//===========================================================================
	// ErrorServiceEmailSendError -> Error sending email: %v.
	ErrorServiceEmailSendError LocalizedString = "error.service.email.send.error"
	// ErrorServiceSmsSendError -> Error sending SMS with code: %v.
	ErrorServiceSmsSendError LocalizedString = "error.service.sms.send.error"
	// ErrorServiceOtpSendError -> Error sending OTP code with error: %v.
	ErrorServiceOtpSendError LocalizedString = "error.service.otp.send.error"

	//===========================================================================
	//  NativeLogin Service
	//===========================================================================
	// ErrorNativeLoginConfigErrors -> NativeLogin service initialized with errors: %+v
	ErrorNativeLoginConfigErrors LocalizedString = "error.native.login.config.errors"

	//===========================================================================
	//  Management API
	//===========================================================================
	// ErrorNativeLoginMaNoKeyID -> No key id found in request to management api.
	ErrorNativeLoginMaNoKeyID LocalizedString = "error.native.login.ma.no.key.id"
	// ErrorNativeLoginMaErrorKeyWithID -> Error getting key with ID: %s, error: %s.
	ErrorNativeLoginMaErrorKeyWithID LocalizedString = "error.native.login.ma.error.key.with.id"
	// ErrorNativeLoginMaErrorSignature -> Invalid signature for request: %s.
	ErrorNativeLoginMaErrorSignature LocalizedString = "error.native.login.ma.error.signature"
	// ErrorNativeLoginMaKeyInactive -> The management key is inactive.
	ErrorNativeLoginMaKeyInactive LocalizedString = "error.native.login.ma.key.inactive"
	// ErrorNativeLoginMaKeyExpired -> The management key is expired.
	ErrorNativeLoginMaKeyExpired LocalizedString = "error.native.login.ma.key.expired"

	//===========================================================================
	//  Admin panel API
	//===========================================================================
	// ErrorAdminPanelNoSkipLimit -> Error parsing Skip and Limit params from request: %s.
	ErrorAdminPanelNoSkipLimit LocalizedString = "error.admin.panel.no.skip.limit"
	// ErrorAdminPanelGetUsers -> Error getting users list: %s.
	ErrorAdminPanelGetUsers LocalizedString = "error.admin.panel.get.users"
	// ErrorAdminPanelNotAuthorized -> You are not authorized to access admin panel.
	ErrorAdminPanelNotAuthorized LocalizedString = "error.admin.panel.not.authorized"
	// ErrorAdminPanelMissingSession -> Session lost or expired, logged out.
	ErrorAdminPanelMissingSession LocalizedString = "error.admin.panel.missing.session"
	// ErrorAdminPanelMissingCookie -> error getting cookie: %v.
	ErrorAdminPanelMissingCookie LocalizedString = "error.admin.panel.missing.cookie"
	// ErrorAdminPanelInvalidCookie -> invalid cookie: %v.
	ErrorAdminPanelInvalidCookie LocalizedString = "error.admin.panel.invalid.cookie"
	// ErrorAdminPanelAdminCredentialsError -> error getting admin credentials: %v
	ErrorAdminPanelAdminCredentialsError LocalizedString = "error.admin.panel.admin.credentials.error"
	// ErrorAdminPanelAdminCredentialsMismatch -> login and password are wrong, admin credentials set via env variables for identifo
	ErrorAdminPanelAdminCredentialsMismatch LocalizedString = "error.admin.panel.admin.credentials.mismatch"
	// ErrorAdminPanelCreateSession -> unable to create session with error: %v
	ErrorAdminPanelCreateSession LocalizedString = "error.admin.panel.create.session"
	// ErrorAdminPanelAdminCredentialsNotSet -> env variables for admin login and password not set
	ErrorAdminPanelAdminCredentialsNotSet LocalizedString = "error.admin.panel.admin.credentials.not.set"
	// ErrorAdminPanelPrivateKeyEncoding -> error encoding private key: %v
	ErrorAdminPanelPrivateKeyEncoding LocalizedString = "error.admin.panel.private.key.encoding"
	// ErrorAdminPanelPrivateKeySave -> error saving private key: %v
	ErrorAdminPanelPrivateKeySave LocalizedString = "error.admin.panel.private.key.save"
	// ErrorAdminPanelPrivateKeyLoad -> error loading private key: %v
	ErrorAdminPanelPrivateKeyLoad LocalizedString = "error.admin.panel.private.key.load"
	// ErrorAdminPanelPrivateKeyEncode -> error encoding private key to PEM: %v
	ErrorAdminPanelPrivateKeyEncode LocalizedString = "error.admin.panel.private.key.encode"
	// ErrorAdminPanelPublicKeyEncode -> error encoding public key to PEM: %v
	ErrorAdminPanelPublicKeyEncode LocalizedString = "error.admin.panel.public.key.encode"
	// ErrorAdminPanelKeyAlgUnsupported -> unsupported algorithm: %v
	ErrorAdminPanelKeyAlgUnsupported LocalizedString = "error.admin.panel.key.alg.unsupported"
	// ErrorAdminPanelGenerateSecret -> error create secret for app: %v
	ErrorAdminPanelGenerateSecret LocalizedString = "error.admin.panel.generate.secret"
	// ErrorAdminPanelAPPSecretLength -> Incorrect app secret string length %d, expecting 24 to 48 symbols inclusively
	ErrorAdminPanelAPPSecretLength LocalizedString = "error.admin.panel.app.secret.length"
	// ErrorAdminPanelAPPSecretNotBase64 -> App secret should be base64 encoded
	ErrorAdminPanelAPPSecretNotBase64 LocalizedString = "error.admin.panel.app.secret.not.base64"

	//===========================================================================
	//  Password
	//===========================================================================
	// PasswordLengthPolicy -> Password length must be at least %d characters.
	PasswordLengthPolicy LocalizedString = "password.length.policy"
	// PasswordRejectCompromised -> Reject compromised passwords, powered by HaveBeenPwned.
	PasswordRejectCompromised LocalizedString = "password.reject.compromised"
	// PasswordRequireLowercase -> Require at least one lowercase characters.
	PasswordRequireLowercase LocalizedString = "password.require.lowercase"
	// PasswordRequireUppercase -> Require at least one uppercase characters.
	PasswordRequireUppercase LocalizedString = "password.require.uppercase"
	// PasswordRequireNumber -> Require at least one number.
	PasswordRequireNumber LocalizedString = "password.require.number"
	// PasswordRequireSymbol -> Require at least special character: !$%%^&*()_+{}:@[];'#<>?,./|\\-=?.
	PasswordRequireSymbol LocalizedString = "password.require.symbol"
)
