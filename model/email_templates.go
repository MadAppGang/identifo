package model

type EmailTemplateType string

const (
	EmailTemplateTypeInvite        EmailTemplateType = "invite-email"
	EmailTemplateTypeResetPassword EmailTemplateType = "reset-password-email"
	EmailTemplateTypeTFAWithCode   EmailTemplateType = "tfa-code-email"
	EmailTemplateTypeVerifyEmail   EmailTemplateType = "verify-email"
	// EmailTemplateTypeWelcome       EmailTemplateType = "welcome-email"

	DefaultTemplateExtension = "html"
)

func (t EmailTemplateType) FileName() string {
	return string(t) + "." + DefaultTemplateExtension
}

func (t EmailTemplateType) String() string {
	return string(t)
}

func AllEmailTemplatesFileNames() []string {
	return []string{
		EmailTemplateTypeInvite.FileName(),
		EmailTemplateTypeResetPassword.FileName(),
		EmailTemplateTypeTFAWithCode.FileName(),
		EmailTemplateTypeVerifyEmail.FileName(),
		// EmailTemplateTypeWelcome.FileName(),
	}
}
