package management

type InvitationTokenRequest struct {
	Email         string                 `json:"email"`
	ApplicationID string                 `json:"application_id"`
	Role          string                 `json:"access_role"`
	CallbackURL   string                 `json:"callback_url"`
	Data          map[string]interface{} `json:"data"`
}

type ResetPasswordTokenRequest struct {
	Email string `json:"email"`
}
