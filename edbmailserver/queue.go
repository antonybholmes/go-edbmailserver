package edbmailserver

const (
	QUEUE_EMAIL_TYPE_VERIFY           = "verify"
	QUEUE_EMAIL_TYPE_VERIFIED         = "verified"
	QUEUE_EMAIL_TYPE_PASSWORDLESS     = "passwordless"
	QUEUE_EMAIL_TYPE_PASSWORD_RESET   = "password-reset"
	QUEUE_EMAIL_TYPE_PASSWORD_UPDATED = "password-updated"
	QUEUE_EMAIL_TYPE_EMAIL_RESET      = "email-reset"
	QUEUE_EMAIL_TYPE_EMAIL_UPDATED    = "email-updated"
	QUEUE_EMAIL_TYPE_ACCOUNT_CREATED  = "account-created"
	QUEUE_EMAIL_TYPE_ACCOUNT_UPDATED  = "account-updated"
	QUEUE_EMAIL_TYPE_OTP              = "otp"
)
